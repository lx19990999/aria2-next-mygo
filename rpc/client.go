package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client talks to an aria2 / aria2-next JSON-RPC endpoint over HTTP.
type Client struct {
	URL    string
	Secret string
	HTTP   *http.Client
}

func New(url, secret string) *Client {
	return &Client{
		URL:    url,
		Secret: secret,
		// aria2-next (and some reverse proxies) often close keep-alive sockets
		// abruptly; that surfaces as intermittent "Post …: EOF" under 1s polling.
		HTTP: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				DisableKeepAlives:   true,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

func isTransientRPCErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "eof") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe")
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc %d: %s", e.Code, e.Message)
}

func (c *Client) call(method string, params ...any) (json.RawMessage, error) {
	args := make([]any, 0, len(params)+1)
	if c.Secret != "" {
		args = append(args, "token:"+c.Secret)
	}
	args = append(args, params...)

	body, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		ID:      "aria2-next-mygo",
		Method:  method,
		Params:  args,
	})
	if err != nil {
		return nil, err
	}

	var last error
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Close = true

		res, err := c.HTTP.Do(req)
		if err != nil {
			last = err
			if attempt == 0 && isTransientRPCErr(err) {
				time.Sleep(40 * time.Millisecond)
				continue
			}
			return nil, err
		}

		raw, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			last = err
			if attempt == 0 && isTransientRPCErr(err) {
				time.Sleep(40 * time.Millisecond)
				continue
			}
			return nil, err
		}
		if res.StatusCode == http.StatusServiceUnavailable || res.StatusCode >= 500 {
			last = fmt.Errorf("HTTP %d", res.StatusCode)
			if attempt == 0 {
				time.Sleep(80 * time.Millisecond)
				continue
			}
			return nil, last
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, fmt.Errorf("HTTP %d", res.StatusCode)
		}
		if len(raw) == 0 {
			last = io.EOF
			if attempt == 0 {
				time.Sleep(40 * time.Millisecond)
				continue
			}
			return nil, last
		}

		var out rpcResponse
		if err := json.Unmarshal(raw, &out); err != nil {
			last = err
			if attempt == 0 && isTransientRPCErr(err) {
				time.Sleep(40 * time.Millisecond)
				continue
			}
			return nil, err
		}
		if out.Error != nil {
			return nil, out.Error
		}
		return out.Result, nil
	}
	return nil, last
}

type Version struct {
	Version          string   `json:"version"`
	Product          string   `json:"product"`
	EnabledFeatures  []string `json:"enabledFeatures"`
	MediaFeatures    []string `json:"mediaFeatures"`
	DownloadFeatures []string `json:"downloadFeatures"`
}

func (c *Client) GetVersion() (Version, error) {
	raw, err := c.call("aria2.getVersion")
	if err != nil {
		return Version{}, err
	}
	var v Version
	return v, json.Unmarshal(raw, &v)
}

type GlobalStat struct {
	DownloadSpeed string `json:"downloadSpeed"`
	UploadSpeed   string `json:"uploadSpeed"`
	NumActive     string `json:"numActive"`
	NumWaiting    string `json:"numWaiting"`
	NumStopped    string `json:"numStopped"`
}

func (c *Client) GetGlobalStat() (GlobalStat, error) {
	raw, err := c.call("aria2.getGlobalStat")
	if err != nil {
		return GlobalStat{}, err
	}
	var s GlobalStat
	return s, json.Unmarshal(raw, &s)
}

// Task is a subset of aria2 tellStatus / tellActive fields.
type Task struct {
	GID             string `json:"gid"`
	Status          string `json:"status"`
	TotalLength     string `json:"totalLength"`
	CompletedLength string `json:"completedLength"`
	DownloadSpeed   string `json:"downloadSpeed"`
	UploadSpeed     string `json:"uploadSpeed"`
	Connections     string `json:"connections"`
	ErrorCode       string `json:"errorCode"`
	ErrorMessage    string `json:"errorMessage"`
	Dir             string `json:"dir"`
	Files           []File `json:"files"`
	BitTorrent      *struct {
		Info *struct {
			Name string `json:"name"`
		} `json:"info"`
	} `json:"bittorrent"`
	Ed2k *struct {
		Name string `json:"name"`
		Hash string `json:"hash"`
	} `json:"ed2k"`
	Media *struct {
		Progress         float64 `json:"progress"`
		DownloadedLength string  `json:"downloadedLength"`
		Protocol         string  `json:"protocol"`
		State            string  `json:"state"`
	} `json:"media"`
}

type File struct {
	Path string   `json:"path"`
	URIs []URIRef `json:"uris"`
}

type URIRef struct {
	URI string `json:"uri"`
}

func (c *Client) tell(method string, offset, num int) ([]Task, error) {
	keys := []string{
		"gid", "status", "totalLength", "completedLength",
		"downloadSpeed", "uploadSpeed", "connections",
		"errorCode", "errorMessage", "dir", "files",
		"bittorrent", "ed2k", "media",
	}
	raw, err := c.call(method, offset, num, keys)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	return tasks, json.Unmarshal(raw, &tasks)
}

func (c *Client) TellActive() ([]Task, error) {
	keys := []string{
		"gid", "status", "totalLength", "completedLength",
		"downloadSpeed", "uploadSpeed", "connections",
		"errorCode", "errorMessage", "dir", "files",
		"bittorrent", "ed2k", "media",
	}
	raw, err := c.call("aria2.tellActive", keys)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	return tasks, json.Unmarshal(raw, &tasks)
}

func (c *Client) TellWaiting(offset, num int) ([]Task, error) {
	return c.tell("aria2.tellWaiting", offset, num)
}

func (c *Client) TellStopped(offset, num int) ([]Task, error) {
	return c.tell("aria2.tellStopped", offset, num)
}

func (c *Client) AddURI(uri string) (string, error) {
	raw, err := c.call("aria2.addUri", []string{uri})
	if err != nil {
		return "", err
	}
	var gid string
	return gid, json.Unmarshal(raw, &gid)
}

func (c *Client) Pause(gid string) error {
	_, err := c.call("aria2.pause", gid)
	return err
}

func (c *Client) Unpause(gid string) error {
	_, err := c.call("aria2.unpause", gid)
	return err
}

func (c *Client) Remove(gid string) error {
	_, err := c.call("aria2.remove", gid)
	return err
}

func (c *Client) ForceRemove(gid string) error {
	_, err := c.call("aria2.forceRemove", gid)
	return err
}

// PurgeDownloadResult clears completed/error/removed results from memory without deleting files.
func (c *Client) PurgeDownloadResult() error {
	_, err := c.call("aria2.purgeDownloadResult")
	return err
}

// RemoveDownloadResult removes one finished result from the list without deleting files.
func (c *Client) RemoveDownloadResult(gid string) error {
	_, err := c.call("aria2.removeDownloadResult", gid)
	return err
}

// Shutdown asks aria2-next to shut down gracefully.
func (c *Client) Shutdown() error {
	_, err := c.call("aria2.shutdown")
	return err
}

// ForceShutdown asks aria2-next to shut down immediately.
func (c *Client) ForceShutdown() error {
	_, err := c.call("aria2.forceShutdown")
	return err
}

// GetGlobalOption returns current global options as key/value strings.
func (c *Client) GetGlobalOption() (map[string]string, error) {
	raw, err := c.call("aria2.getGlobalOption")
	if err != nil {
		return nil, err
	}
	var opts map[string]string
	if err := json.Unmarshal(raw, &opts); err != nil {
		return nil, err
	}
	return opts, nil
}

// ChangeGlobalOption updates one or more global options.
func (c *Client) ChangeGlobalOption(opts map[string]string) error {
	if len(opts) == 0 {
		return nil
	}
	_, err := c.call("aria2.changeGlobalOption", opts)
	return err
}

func (t Task) Name() string {
	if t.BitTorrent != nil && t.BitTorrent.Info != nil && t.BitTorrent.Info.Name != "" {
		return t.BitTorrent.Info.Name
	}
	if t.Ed2k != nil && t.Ed2k.Name != "" {
		return t.Ed2k.Name
	}
	if len(t.Files) > 0 {
		if p := t.Files[0].Path; p != "" {
			if i := strings.LastIndexAny(p, `/\`); i >= 0 && i+1 < len(p) {
				return p[i+1:]
			}
			return p
		}
		if len(t.Files[0].URIs) > 0 {
			return t.Files[0].URIs[0].URI
		}
	}
	return t.GID
}

func (t Task) Progress() float64 {
	if t.Media != nil && t.Media.Progress > 0 {
		if t.Media.Progress <= 1 {
			return t.Media.Progress
		}
		return t.Media.Progress / 100
	}
	total := ParseInt(t.TotalLength)
	if total <= 0 {
		if t.Status == "complete" {
			return 1
		}
		return 0
	}
	return float64(ParseInt(t.CompletedLength)) / float64(total)
}

func ParseInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func FormatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func FormatSpeed(s string) string {
	return FormatBytes(ParseInt(s)) + "/s"
}

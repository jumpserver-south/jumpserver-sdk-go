package audits

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/sdkutil"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

// FTP log URL constants.
const (
	FTPLogListURL   = "/api/v1/audits/ftp-logs/"
	FTPLogDetailURL = "/api/v1/audits/ftp-logs/%s/"
	FTPLogUploadURL = "/api/v1/audits/ftp-logs/%s/upload/"
)

// ListFTPLogs returns a paginated list of FTP logs.
func (s *Service) ListFTPLogs(ctx context.Context, opts *core.ListOptions) ([]model.FTPLog, *core.Response, error) {
	return sdkutil.List[model.FTPLog](ctx, s.client, FTPLogListURL, opts)
}

// UploadFTPFile uploads a file associated with an FTP log entry.
func (s *Service) UploadFTPFile(ctx context.Context, ftpLogID, filePath string) (*core.Response, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	mpw := multipart.NewWriter(body)
	part, err := mpw.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	if err := mpw.Close(); err != nil {
		return nil, err
	}

	url := sdkutil.Spath(FTPLogUploadURL, ftpLogID)
	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", mpw.FormDataContentType())
	httpReq.Body = io.NopCloser(body)
	httpReq.ContentLength = int64(body.Len())
	return s.client.Do(ctx, httpReq, nil)
}

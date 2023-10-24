package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gnames/gndiff/internal/io/ingestio"
	gndiff "github.com/gnames/gndiff/pkg"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/output"
	"github.com/gnames/gndiff/pkg/ent/record"
	echo "github.com/labstack/echo/v4"
)

func apiInfo(c echo.Context) error {
	return c.String(http.StatusOK,
		`The API is described at
https://apidoc.globalnames.org/gndiff`)
}

func apiDiffPOST(gnd gndiff.GNdiff) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := getContext(c)
		defer cancel()
		chErr := make(chan error)

		go func() {
			defer close(chErr)

			var err error
			var qry, ref []record.Record
			var res output.Output
			var params Input

			err = c.Bind(&params)

			if err == nil {
				slog.Info("match",
					slog.String("method", "POST"),
				)
			}

			qryFileName := params.QueryFileName
			if strings.TrimSpace(qryFileName) == "" {
				qryFileName = "query"
			}
			refFileName := params.ReferenceFileName
			if strings.TrimSpace(refFileName) == "" {
				refFileName = "reference"
			}

			qry, err = textToRecords(
				params.QueryText,
				qryFileName,
			)

			if err == nil {
				ref, err = textToRecords(
					params.ReferenceText,
					refFileName,
				)
			}
			if err == nil {
				res, err = gnd.Compare(qry, ref)
			}

			if err == nil {
				err = c.JSON(http.StatusOK, res)
			}

			chErr <- err
		}()

		select {
		case <-ctx.Done():
			<-chErr
			return ctx.Err()
		case err := <-chErr:
			return err
		case <-time.After(6 * time.Minute):
			return errors.New("request took too long")
		}
	}
}

func textToRecords(s, fileName string) ([]record.Record, error) {
	cfg := config.New()
	ing := ingestio.New(cfg)
	ext := strings.ToLower(filepath.Ext(fileName))
	r := strings.NewReader(s)
	records, sep, err := ing.RecordsFromText(r, fileName, ext)
	if err != nil {
		return nil, err
	}
	if records != nil {
		return records, nil
	}

	// recreate reader insteaf of rewinding
	r = strings.NewReader(s)
	return ing.RecordsFromCSV(r, sep, fileName)
}

func getContext(c echo.Context) (ctx context.Context, cancel func()) {
	ctx = c.Request().Context()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	return ctx, cancel
}

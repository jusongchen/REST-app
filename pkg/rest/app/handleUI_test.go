package app

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStaticHTMLHandler(t *testing.T) {
	t.Parallel()
	t.Run("Happy_path ", func(t *testing.T) {

		urlPath := "/UI"
		staticFilePath := "./testdata/UI"
		for _, allowCORS := range []bool{false, true} {

			h := StaticHTMLHandler(urlPath, staticFilePath, allowCORS)
			require.HTTPSuccess(t, h, "GET", urlPath, nil)
			require.HTTPBodyContains(t, h, "GET", urlPath, nil, "Welcome!")
		}

	})

	t.Run("html_file_dir_not_exists", func(t *testing.T) {

		urlPath := "/UI/"
		staticFilePath := "./testdata/dir-not-exists"
		for _, allowCORS := range []bool{false, true} {

			h := StaticHTMLHandler(urlPath, staticFilePath, allowCORS)

			require.HTTPStatusCode(t, h, "GET", urlPath, url.Values{}, 500)
			require.HTTPBodyContains(t, h, "GET", urlPath, nil, "not exist")

		}

	})
}

# JSONRPC version 2 over Http

Gets code from https://github.com/powerman/rpc-codec
Only modify http.go to remove strict HTTP content type checking


    else if resp.Header.Get("Content-Type") != contentType {
            err = fmt.Errorf("bad HTTP Content-Type: %s", resp.Header.Get("Content-Type"))
        }


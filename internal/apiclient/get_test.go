package apiclient

// TestCallAPISuccess

//     httptest.NewServer returns status 200 and body []byte("ok")
//     Assert returned bytes equal "ok" and error is nil

// TestCallAPIBadStatus

//     httptest.NewServer returns status 404
//     Assert error is non-nil

// TestCallAPIBadURL

//     Call CallAPI("http://127.0.0.1:1") (nothing listening)
//     Assert error is non-nil

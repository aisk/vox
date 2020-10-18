---
title: Response
nav_order: 6
---

# Response

A response is what to send to the user's browser.

## Response.Status

This is straightforward that you can set this field, to change the response's status code. The default value of this field is 200 if this request has been processed by any routers. Otherwise, it should be 404, even though it hits some middlewares.

## Response.Body

This field is a bit of magic since it's type is `interface{}`, which means that you can put any type of values to this field. But different types will have different behaviors to marshal its value to the user's browser.

1) If the `Response.Body` is a `[]bytes` or `string`, the value will be written to the user's browser.
1) If the `Response.Body` is an `io.Reader`, vox will read all the data and write them to the user's browser.
1) If the `Response.Body` is an `io.ReadCloser`, vox will do the same behavior as above, plus close it.
1) Otherwise, vox will try to respect request header's `Accept` filed to marshal `Response.Body` to browser, with JSON as a fallback. Unfortunately, only JSON is supported now. Pull Requests are welcome.

## Response.Header

This field is the `http.Header` which will be written to the browser, you can set it directly as what you do when you're using the native HTTP library.

The response header's `Content-Type` field is an exception. It will be set with the most suitable value with the `Response.Body` if you didn't set it explicitly. For example, if the response body is HTML, it's value will be set to `text/html`.

## Response.Writer

This is the hidden secret of vox. If the current mechanism is not suitable for your scenario, you can use it as the underlying `http.ResponseWriter`. But before this, you should set the `Response.DontRespond` to `false`.
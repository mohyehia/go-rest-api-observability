# go-rest-api-observability
Sample GO Rest API along with observability setup using prometheus, grafana, tempo &amp; loki


## Table Driven Tests
The standard, most highly recommended pattern in the Go community to solve this is Table-Driven Tests.
It allows you to define a list (slice) of test cases and run them inside a single loop using `t.Run()`.

Example:
```go
func TestPostsClient_GetPost(t *testing.T) {
    // 1. Define the data structure for our test cases
    type testCase struct {
        name           string
        postID         int
        mockStatus     int
        mockResponse   *Post
        expectedResult *Post
        expectedErr    string
    }

    // 2. Define the "Table" of scenarios
    tests := []testCase{
        {
            name:       "Success Scenario",
            postID:     1,
            mockStatus: http.StatusOK,
            mockResponse: &Post{
                UserID: 1, ID: 1, Title: "Test Post", Body: "Body content",
            },
            expectedResult: &Post{
                UserID: 1, ID: 1, Title: "Test Post", Body: "Body content",
            },
            expectedErr: "",
        },
        {
            name:           "Bad Request Error Scenario",
            postID:         1,
            mockStatus:     http.StatusBadRequest,
            mockResponse:   nil,
            expectedResult: nil,
            expectedErr:    "status code 400",
        },
    }

    // 3. Iterate over the matrix table execution
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                expectedPath := fmt.Sprintf("/posts/%d", tc.postID)
                if req.URL.Path != expectedPath {
                    t.Errorf("Expected path %s, got %s", expectedPath, req.URL.Path)
                }
                
                w.WriteHeader(tc.mockStatus)
                if tc.mockResponse != nil {
                    w.Header().Set("Content-Type", "application/json")
                    _ = json.NewEncoder(w).Encode(tc.mockResponse)
                }
            }))
            defer server.Close()

            pc := NewPostsClient(http.DefaultClient, server.URL)
            post, err := pc.GetPost(tc.postID)

            // Check Error Assertions
            if tc.expectedErr != "" {
                if err == nil {
                    t.Fatalf("Expected error containing '%s', got nil", tc.expectedErr)
                }
                if err.Error() != tc.expectedErr {
                    t.Errorf("Expected error message '%s', got '%s'", tc.expectedErr, err.Error())
                }
                return
            }

            if err != nil {
                t.Fatalf("Unexpected error: %v", err)
            }

            // Clean Direct Assertions instead of reflect.DeepEqual
            if post.ID != tc.expectedResult.ID || post.Title != tc.expectedResult.Title {
                t.Errorf("Expected post %+v, got %+v", tc.expectedResult, post)
            }
        })
    }
}
```
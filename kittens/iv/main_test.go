package main

import (
    "io/ioutil"
    "os"
    "testing"
    "fmt"
    "path/filepath"
)

func TestIsImage(t *testing.T) {
    cases := []struct {
        filename string
        expected bool
    }{
        {"image.jpg", true},
        {"image.jpeg", true},
        {"image.png", true},
        {"image.gif", true},
        {"image.bmp", true},
        {"image.txt", false},
        {"image", false},
    }

    for _, c := range cases {
        result := isImage(c.filename)
        if result != c.expected {
            t.Errorf("isImage(%q) == %v, expected %v", c.filename, result, c.expected)
        }
    }
}

func TestDiscoverImages(t *testing.T) {
    // Create a temporary directory
    tempDir, err := ioutil.TempDir("", "testDiscoverImages")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    // Create some test files
    files := []string{"image.jpg", "image.png", "text.txt"}
    for _, file := range files {
        _, err := os.Create(filepath.Join(tempDir, file))
        if err != nil {
            t.Fatal(err)
        }
    }

    // Discover images
    err = discoverImages(tempDir)
    if err != nil {
        t.Fatalf("discoverImages returned an error: %v", err)
    }

    // Check the results
    expectedImages := []string{filepath.Join(tempDir, "image.jpg"), filepath.Join(tempDir, "image.png")}
    if len(globalImages) != len(expectedImages) {
        t.Errorf("expected %d images, got %d", len(expectedImages), len(globalImages))
    }

    for i, img := range globalImages {
        if img != expectedImages[i] {
            t.Errorf("expected %s, got %s", expectedImages[i], img)
        }
    }
}

func TestLoadConfig(t *testing.T) {
    // Create a temporary config file
    configContent := `
windowParam:
  xParam: 4
  yParam: 5
`
    configFile, err := ioutil.TempFile("", "config.yaml")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(configFile.Name())

    if _, err := configFile.Write([]byte(configContent)); err != nil {
        t.Fatal(err)
    }
    if err := configFile.Close(); err != nil {
        t.Fatal(err)
    }

    // Load the config
    err = loadConfig(configFile.Name())
    if err != nil {
        t.Fatalf("loadConfig returned an error: %v", err)
    }

    // Check the results
    fmt.Printf("globalConfig after loadConfig: %+v\n", globalConfig) // Debug output

    if globalConfig.gridParam.x_param != 4 {
        t.Errorf("expected x_param to be 4, got %d", globalConfig.gridParam.x_param)
    }
    if globalConfig.gridParam.y_param != 5 {
        t.Errorf("expected y_param to be 5, got %d", globalConfig.gridParam.y_param)
    }
}

func TestGetWindowSize(t *testing.T) {
    // Note: Testing the actual window size might not be feasible in a unit test environment
    // because it depends on the terminal environment. This is just a basic structure for the test.
    err := getWindowSize(globalWindowParameters)
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }

    fmt.Printf("globalWindowParameters after getWindowSize: %+v\n", globalWindowParameters) // Debug output

    if globalWindowParameters.Row == 0 || globalWindowParameters.Col == 0 {
        t.Errorf("expected window size to be set, got rows: %d, cols: %d", globalWindowParameters.Row, globalWindowParameters.Col)
    }
}

func TestPaginateImages(t *testing.T) {
    // Set the global images
    globalImages = []string{"img1.jpg", "img2.jpg", "img3.jpg", "img4.jpg", "img5.jpg", "img6.jpg"}

    // Set the grid parameters
    globalConfig.gridParam = gridConfig{x_param: 2, y_param: 2}

    // Paginate the images
    paginateImages()

    // Check the results
    expectedPages := [][]string{
        {"img1.jpg", "img2.jpg"},
        {"img3.jpg", "img4.jpg"},
    }

    if len(globalImagePages) != len(expectedPages) {
        t.Errorf("expected %d pages, got %d", len(expectedPages), len(globalImagePages))
    }

    for i, page := range globalImagePages {
        for j, img := range page {
            if img != expectedPages[i][j] {
                t.Errorf("expected %s, got %s", expectedPages[i][j], img)
            }
        }
    }
}

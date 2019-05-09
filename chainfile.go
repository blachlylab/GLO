package glo

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "net/http"
    "compress/gzip"
)


type ChainFile struct {
    referenceBuild, queryBuild, filepath string
    chains map[string][]*Chain
}

func (cf *ChainFile) add(chain *Chain) {
    contig := chain.tName
    cf.chains[contig] = append(cf.chains[contig], chain)
}

// fileExists checks if a file exists at the specified filepath, 
// returning a corresponding boolean value.
func fileExists(filepath string) bool {
    file, err := os.Open(filepath)
    file.Close()
    if err == nil {
        // File existed
        return true
    }
    return false
}

// getFileType deterimines the MIME file type of the file located
// at the specified filepath, by inspecting the first 512 bytes of
// the file using http.DetectContentType(). The file type is returned
// as a string.
func getFileType(filepath string) string {
    if !fileExists(filepath) {
        fmt.Printf("Error: Specified file (%s) cannot be opened.\n", filepath)
        os.Exit(1)
    }

    file, _ := os.Open(filepath)
    defer file.Close()

    // Read the first 512 bytes from the file into the buffer. This should
    // be sufficient to determine the filetype.
    buffer := make([]byte, 512)
    _, err := file.Read(buffer)

    if err != nil {
        fmt.Printf("Error opening file: %s\n", err)
        os.Exit(1)
    }
    file.Close()

    return string(http.DetectContentType(buffer))
}

func get_chunk(reader *bufio.Reader, i int) string {
    s := ""
    p, p_err := reader.Peek(6)
    if p_err == nil {
        s = string(p)
        pos := strings.Index(s, "\n")
        if pos != -1 {
            s = s[:pos]
        }
    }
    return(string(s))
}

// ChainFile.Load loads the contents of the specified chain file
// into the struct, storing them as a map of chains.
func (cf *ChainFile) Load() {
    cf.chains = make(map[string][]*Chain)

    // Determine if this is a plain text (uncompressed) chain file
    // or a gzipped file. No other formats are supported at this time.
    filetype := getFileType(cf.filepath)
    if filetype[:11] == "text/plain;" {
        filetype = "text"
    } else if filetype == "application/x-gzip" {
        filetype = "gzip"
        //os.Exit(1)
    } else {
        fmt.Printf("Unsupported filetype: %s\n", filetype)
        os.Exit(1)
    }

    // getFileType already checks if the file exists, so we don't
    // need to repeat that call here.
    f, _ := os.Open(cf.filepath)
    defer f.Close()

    // Conditionally either create a gzip or plain text reader
    var reader *bufio.Reader
    if filetype == "gzip" {
        fgz, _ := gzip.NewReader(f)
        reader = bufio.NewReader(fgz)
    } else {
        reader = bufio.NewReader(f)
    }

    // Use the reader to can through the entire file one
    // line at a time, building the ChainFile structure full
    // of Chain structs.
    var chain *Chain
    for {
        // Peek ahead one byte to determine if the Reader can
        // read any further data.
        _, p_err := reader.Peek(1)
        if p_err != nil {
            // Peek encountered an error; assuming it is EOF
            // so that the for look needs to terminate.
            break
        }

        // Use helper function to peek at the up to next 6 bytes,
        // returned as a string that contains anything before the
        // newline character.
        s := strings.TrimSpace(get_chunk(reader, 6))

        if len(s) == 0 {
            // Blank line, read until new line and discard.
            reader.ReadString('\n')
        } else {
            if len(s) > 4 && s[:5] == "chain" {
                // Start of a new chain.
                line, _ := reader.ReadString('\n')
                if chain != nil {
                    // Store this current chain before
                    // initializing a new one.
                    cf.add(chain)
                }
                chain = new(Chain)
                chain.FromString(string(line))
            } else if s[0] != '#' {
                // This should be mapping block(s) that
                // will get loaded into the current chain.
                chain.load_links(reader)
            } else {
                // Junk line, should never hit this in a
                // well-formatted chain file.
                reader.ReadString('\n')
            } // end inner if-elseif-else block
        } // end outer if-else block
    } // end for loop

    if chain != nil {
        // Something unstored remains in a chain; store it.
        cf.add(chain)
        chain = nil
    }
}

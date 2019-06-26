package common

// This file wraps zip functions

type FileEntry struct {
	SourceName, ZipName string
}

type Files []*FileEntry

//Zip wraps zip functionality

func Zip(outputFile string, files Files) error {
	// Create a buffer to write our archive to.
	/*	buf := new(bytes.Buffer)

		// Create a new zip archive.
		w := zip.NewWriter(buf)

		for _, file := range files {
			f, err := w.Create(file.SourceName)
			if err != nil {
				return err
			}

					_, err = f.Write([]byte(file.Body))
			if err != nil {
				return err
			}
		}

		// Make sure to check the error on Close.
		err := w.Close()
		if err != nil {
			return err
		}
	*/
	return nil
}

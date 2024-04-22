package fileindexing

import (
	"encoding/json"
	"os"
	"slices"

	t "github.com/JackPairce/MicroService/services/types"
)

type FileIndexing map[string][]*t.File

func SaveData(data FileIndexing, path string) error {
	var pdData JSONData
	for key, value := range data {
		keyValue := &KeyValue{
			Key:    key,
			Values: value,
		}
		pdData.Data = append(pdData.Data, keyValue)
	}
	serializedData, err := json.Marshal(&pdData)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(serializedData)
	return err
}

func LoadData(path string) (FileIndexing, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pdData JSONData
	if err := json.Unmarshal(file, &pdData); err != nil {
		return nil, err
	}
	jsonData := make(FileIndexing)
	for _, value := range pdData.Data {
		jsonData[value.Key] = value.Values
	}
	return jsonData, nil
}
func generateKeys(input string) []string {
	var keys []string

	var generate func(prefix, remaining string)
	generate = func(prefix, remaining string) {
		if len(remaining) == 0 {
			return
		}
		keys = append(keys, prefix)
		for i := 0; i < len(remaining); i++ {
			generate(prefix+string(remaining[i]), remaining[i+1:])
		}
	}

	// generate("", strings.Split(input, ".")[0])
	generate("", input)

	return keys
}

func (idx FileIndexing) GetUniqueFileNames(id int32) []*t.File {
	var uniqueFiles []*t.File
	seen := make(map[string]bool)

	for _, files := range idx {
		for _, file := range files {
			if !slices.Contains(file.Ownerid, id) && !seen[file.Name] {
				seen[file.Name] = true
				uniqueFiles = append(uniqueFiles, file)
			}
		}
	}

	return uniqueFiles
}

func (idx FileIndexing) AddFile(f *t.File) {
	// Generate indexes from fileName
	prefixes := generateKeys(f.Name)
	// Add the file name to the index for each prefix
	for _, prefix := range prefixes {
		idx[prefix] = append(idx[prefix], f)
	}
}

func (idx FileIndexing) SearchFiles(queryPrefix string) []*t.File {
	return idx[queryPrefix]
}

func (idx FileIndexing) UpdateFile(File *t.File) {
	prefixes := generateKeys(File.Name)
	for _, prefix := range prefixes {
		files := idx[prefix]
		for i, file := range files {
			files[i] = &t.File{
				Name:    file.Name,
				Path:    file.Path,
				Ownerid: append(file.Ownerid, File.Ownerid...),
				Size:    file.Size,
			}
			idx[prefix] = files
		}
	}
}

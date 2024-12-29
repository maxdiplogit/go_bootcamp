package maps

import "fmt"

func main() {
	websites := map[string]string{
		"Google": "www.google.com",
		"AWS":    "www.aws.com",
	}

	fmt.Println(websites["Google"])
}

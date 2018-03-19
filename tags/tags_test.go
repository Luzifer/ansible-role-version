package tags

import "log"

func ExampleGetLatestTag() {
	tag, err := GetLatestTag("https://github.com/luzifer-ansible/deploy-git.git", true)
	if err != nil {
		log.Fatalf("Could not resolve latest tag: %s", err)
	}

	log.Printf("Latest tag: %q, created at %v", tag.Name, tag.When)
}

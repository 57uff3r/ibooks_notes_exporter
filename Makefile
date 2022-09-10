release_test:
	goreleaser release --snapshot --skip-publish --rm-dist

release_publish:
	goreleaser release --rm-dist
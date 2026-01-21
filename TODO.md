# TODO: Fix Certificate Issues

## Tasks to Complete
- [x] Create `internal/output/certificate.go` with `PrintCertificate` function that accepts `*analyzer.CertificateData` and prints human-readable certificate info using existing output styles.
- [x] Update `GetCommits` in `internal/github/commits.go` to paginate like `GetContributors`: use `per_page=100`, loop over pages, append commits, stop when a page has fewer than `per_page` commits.
- [x] Build the project to verify no compilation errors.
- [x] Test the certificate command to ensure proper output.

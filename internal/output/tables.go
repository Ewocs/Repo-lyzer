package output

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/github"

	"github.com/olekukonko/tablewriter"
)

func PrintRepo(r *github.Repo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Repository", "Stars", "Forks", "Open Issues"})
	table.Append([]string{
		r.FullName,
		fmt.Sprint(r.Stars),
		fmt.Sprint(r.Forks),
		fmt.Sprint(r.OpenIssues),
	})

	table.Render()
}

func PrintContributors(contributors []github.Contributor) {
	if len(contributors) == 0 {
		fmt.Println("No contributors found.")
		return
	}

	fmt.Println("\n👥 Top Contributors")

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Contributor", "Commits", "Avatar"})

	// Show top 10 contributors or all if less than 10
	maxContributors := 10
	if len(contributors) < maxContributors {
		maxContributors = len(contributors)
	}

	for i := 0; i < maxContributors; i++ {
		contributor := contributors[i]
		table.Append([]string{
			contributor.Login,
			fmt.Sprintf("%d", contributor.Commits),
			contributor.AvatarURL,
		})
	}

	table.Render()

	if len(contributors) > 10 {
		fmt.Printf("... and %d more contributors\n", len(contributors)-10)
	}
}

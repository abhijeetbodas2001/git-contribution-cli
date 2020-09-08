## Git-contribution-cli

This is a CLI to generate contribution graph similar to GitHub profile contribution graph, but for local git repositories.

To run the program:
1. Install go lang
2. Install dependencies with `go install`
3. Add the folder paths, in subfolders(any depth) of which your git repos are present, using `go run . -add "<directory-here>"`.  
This will generate a `.gogitlocalstats` file in you home directory, containing paths to all local git repositories.
4. `go run . -email "<emailid-here>"` to generate the graph.

Note: The Windows command console doesn't support output coloring by default. Run this in Windows Terminal or Powershell
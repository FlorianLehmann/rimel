import git

def main():
    repo = git.Repo('.')
    for commit in repo.iter_commits('master'):        
        # do something with the commit, for example:
        print(commit)
        # commit.committed_datetime.strftime("%a %d. %b %Y")
        # commit.message.rstrip() 

        # see http://gitpython.readthedocs.io/en/stable/tutorial.html#the-commit-object for more information on the available attributes/methods of the commit object
if __name__ == "__main__":
    main()
// Fetch commits from the API and display them
async function fetchCommits() {
    try {
        const response = await fetch('/api/commits');
        const commits = await response.json();

        const commitList = document.getElementById('commit-list');
        commitList.innerHTML = ''; // Clear previous commits

        commits.forEach(commit => {
            const commitItem = document.createElement('div');
            commitItem.classList.add('commit');

            commitItem.innerHTML = `
                <h3>Commit: ${commit.sha}</h3>
                <p>Message: ${commit.commit.message}</p>
                <p>Author: ${commit.commit.author.name} (${commit.commit.author.email})</p>
            `;

            commitList.appendChild(commitItem);
        });
    } catch (error) {
        console.error('Error fetching commits:', error);
    }
}

// Fetch commits on page load
window.onload = fetchCommits;
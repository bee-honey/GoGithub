// Fetch commits from the API and display them
async function fetchCommits() {
    try {
        const response = await fetch('/api/commits');
        const commits = await response.json();

        console.log(commits);

        const commitList = document.getElementById('commit-list');
        commitList.innerHTML = ''; // Clear previous commits

        commits.forEach((commit, index) => {
            const commitItem = document.createElement('div');
            commitItem.classList.add('commit');

            // Determine the color based on the conclusion
    let conclusionColor;
    switch (commit.conclusion.toLowerCase()) {
        case 'success':
            conclusionColor = 'green';
            break;
        case 'failure':
            conclusionColor = 'red';
            break;
        default:
            conclusionColor = 'grey';  // For cases where there is no conclusion or 'nothing'
            break;
    }

            // Add commit details
            commitItem.innerHTML = `
                <h3>Commit: ${commit.sha}</h3>
                <p>Message: ${commit.message}</p>
                <p>Author: ${commit.author}</p>
                <p>Status: ${commit.status} (<span style="color: ${conclusionColor};">${commit.conclusion}</span>)</p>
                
            `;

            // Create the Release button
            const releaseButton = document.createElement('button');
            releaseButton.textContent = 'Release';
            releaseButton.id = `release-btn-${index}`;
            releaseButton.disabled = commit.conclusion.toLowerCase() !== 'success';
            releaseButton.addEventListener('click', async () => {
                
                try {
                    const response = await fetch('/api/release', {
                        method: 'POST',
                    });
                    if (response.ok) {
                        alert(`CI-INFO: Workflow triggered successfully for commit: ${commit.sha}`);
                    } else {
                        alert(`CI-ERROR: Failed to trigger workflow for commit: ${commit.sha}`);
                    }
                } catch (error) {
                    console.error('CI-ERROR: Error triggering the release', error)
                }
            });

            // Append the button to the commit item
            commitItem.appendChild(releaseButton);

            // Append the commit item to the commit list
            commitList.appendChild(commitItem);
        });
    } catch (error) {
        console.error('Error fetching commits:', error);
    }
}

// Fetch commits on page load
window.onload = fetchCommits;
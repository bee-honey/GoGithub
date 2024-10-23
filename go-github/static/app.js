// Fetch commits from the API and display them
async function fetchCommits() {
    try {
        const response = await fetch('/api/commits');
        const commits = await response.json();

        const commitList = document.getElementById('commit-list');
        commitList.innerHTML = '';
        commits.forEach(commit => {
            const commitItem = document.createElement('div');
            commitItem.classList.add('commit');

            commitItem.innerHTML = `
                <h3>Commit: ${commit.sha}</h3>
                <p>Message: ${commit.commit.message}</p>
                <p>Author: ${commit.commit.author.name} (${commit.commit.author.email})</p>
            `;

            //Release button
            const releaseButton = document.createElement('button');
            releaseButton.textContent = 'Release';
            releaseButton.id = `release-btn-${index}`;
            releaseButton.addEventListener('click', () => {
                alert(`Releasing commit: ${commit.sha}`);
                // Need to add more logic here to handle the release process
            });
            commitItem.appendChild(releaseButton);
            commitList.appendChild(commitItem);
        });
    } catch (error) {
        console.error('Error fetching commits:', error);
    }
}

// Fetch commits on page load
window.onload = fetchCommits;
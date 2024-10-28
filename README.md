# GoGithub

GoGithub is a simple web-based dashboard developed using Go and JavaScript, designed to monitor the status of GitHub repositories. With GoGithub, you can view recent commits and their statuses, check workflow runs, and even trigger GitHub Actions workflows directly from the dashboard.

## Features
- **Commit and Workflow Overview**: View recent commits, their authors, messages, and statuses.
- **Trigger GitHub Actions**: Use a single button to initiate workflows, like `release.yml`, directly from the dashboard.
- **Authentication**: Uses GitHub Personal Access Tokens for authenticated requests, allowing increased rate limits and secure POST requests.

## Setup

### Prerequisites
- **Go**: Ensure Go is installed and set up on your machine.
- **JavaScript**: Basic frontend JavaScript for UI functionality.
- **GitHub Personal Access Token**: Create a personal access token with `repo` and `workflow` permissions to enable secure API access.

### Parameters
The following environment variables are required to run the application:

| Parameter          | Description                               | Example                |
|--------------------|-------------------------------------------|------------------------|
| `repoOwner`        | Owner of the GitHub repository            | `"bee-honey"`          |
| `repoName`         | Name of the GitHub repository             | `"CICDAbt"`            |
| `workflowFile`     | GitHub Actions workflow file (optional)   | `"release.yml"`        |
| `GITHUB_TOKEN`     | GitHub Personal Access Token              | `your_personal_token`  |

### Installation and Running the Application

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/GoGithub.git
   cd GoGithub


#Output

<img width="1910" alt="Screenshot 2024-10-28 at 11 26 10 AM" src="https://github.com/user-attachments/assets/0023bb85-ee8a-4a23-ab17-7e1504e44294">


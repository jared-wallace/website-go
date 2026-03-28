---
phase: quick
plan: 260327-vgp
type: execute
wave: 1
depends_on: []
files_modified: [README.md]
autonomous: true
---

<objective>
Add a "Production Deployment" section to README.md that documents how to deploy website-go on the existing AWS infrastructure (ALB + Nginx + EC2/ASG + EBS volume defined in the sibling aws-infra repo).

Purpose: A developer (including future-you) can read the README and deploy a new version of the blog to production without archaeology across two repos.
Output: Updated README.md with a Production Deployment section.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@README.md
@CLAUDE.md
@Makefile
@docker-compose.dev.yml

The aws-infra repo (../aws-infra) defines the production environment:

**Infrastructure topology (from main.tf + user_data.sh):**
- Region: us-east-1
- VPC with dual-stack (IPv4 + IPv6) public subnets across 2 AZs
- ALB (main-alb): internet-facing, dualstack, TLS termination via ACM wildcard cert for jared-wallace.com
- ALB listener: port 80 redirects to 443; port 443 forwards to target group on port 80
- ASG (web-asg): min=1, max=2, desired=1; launches t4g.micro (ARM64) Amazon Linux 2023 instances
- EC2 security group: allows port 80 from ALB only, port 22 from anywhere
- EBS volume (gp3, 10GB): attached via user_data.sh to /dev/nvme1n1, mounted at /var/www/html
- user_data.sh installs Docker, AWS CLI, git, and Nginx; creates /var/www/html/app/ for application code
- Nginx reverse proxy: listens on port 80, proxies to localhost:8080 (the Go app's Docker container port)
- systemd service (jw-blog.service): runs /var/www/html/start-app.sh on boot
- Route53: A/AAAA records for jared-wallace.com, www, admin — all aliased to the ALB
- Lambda + EventBridge: updates ssh.jared-wallace.com DNS when ASG launches/terminates instances

**Key paths on the EC2 instance:**
- /var/www/html/ — EBS mount point (persistent across instance replacement)
- /var/www/html/app/ — Application code directory
- /var/www/html/start-app.sh — Deployment entry point (calls deploy.sh if present)

**What does NOT exist yet:**
- No Dockerfile in website-go (Phase 6 work, not yet built)
- No deploy.sh script
- No docker-compose.prod.yml
- No production DATABASE_URL or env var documentation

The README section should document the architecture and deployment steps accurately based on what the infra provides, while noting that the Dockerfile and deploy.sh are Phase 6 deliverables not yet created.
</context>

<tasks>

<task type="auto">
  <name>Task 1: Add Production Deployment section to README.md</name>
  <files>README.md</files>
  <action>
Insert a new "## Production Deployment" section in README.md between "## Continuous Integration" and "## License". The section should cover:

1. **Architecture Overview** — A concise text description of the request flow:
   - Client -> Route53 (jared-wallace.com) -> ALB (TLS termination) -> EC2 (t4g.micro ARM64) -> Nginx (:80) -> Go app (:8080)
   - Persistent storage on EBS volume mounted at /var/www/html
   - Single-instance ASG for self-healing (not horizontal scaling)
   - Note that infrastructure is managed in a separate `aws-infra` repo via Terraform

2. **Deployment Steps** — What a developer actually does to ship a new version. Since Dockerfile/deploy.sh don't exist yet, frame this as "once Phase 6 is complete, the deployment flow will be":
   - SSH to the instance: `ssh ec2-user@ssh.jared-wallace.com`
   - Pull latest code into /var/www/html/app/
   - Build the Docker image (ARM64 target): `docker build -t website-go:latest .`
   - Run via docker-compose (production config) which starts the Go app on :8080 and Postgres with data on the EBS volume
   - Nginx (already running via systemd) proxies traffic to the container

3. **Persistent Storage** — Document the EBS volume layout:
   - /var/www/html/ — EBS mount (survives instance replacement)
   - /var/www/html/app/ — Application code
   - /var/www/html/postgres-data/ — Postgres data directory (planned, per CLAUDE.md stack notes)
   - /var/www/html/images/ — Uploaded images (planned)

4. **Important Caveats** — Pull from STATE.md blockers:
   - ASG max_size must be 1 and EBS delete_on_termination=false before writing production data
   - Postgres EBS bind-mount requires `chown 999:999` on the postgres-data directory before first run
   - Note: Dockerfile, deploy.sh, and docker-compose.prod.yml are Phase 6 deliverables (not yet created)

Keep the tone consistent with the existing README — terse, scannable, code-block-heavy. No fluff. Use the same markdown heading hierarchy (## for section, ### for subsections). Do NOT use emojis.
  </action>
  <verify>
    <automated>grep -c "Production Deployment" README.md | grep -q "1" && grep -q "ssh.jared-wallace.com" README.md && grep -q "EBS" README.md && echo "PASS" || echo "FAIL"</automated>
  </verify>
  <done>README.md contains a Production Deployment section with architecture overview, deployment steps, persistent storage layout, and caveats. The section is factually accurate against the aws-infra Terraform configuration.</done>
</task>

</tasks>

<verification>
- README.md has exactly one "## Production Deployment" heading
- Section appears between CI and License sections
- Architecture description matches aws-infra topology (ALB, Nginx, port 8080, EBS at /var/www/html)
- No claims about files that don't exist yet (Dockerfile, deploy.sh) without noting they are Phase 6 work
- Markdown renders cleanly (no broken formatting)
</verification>

<success_criteria>
A developer reading the Production Deployment section understands the full request flow, knows how to SSH in and deploy, knows where persistent data lives, and is warned about the EBS/ASG caveats — all without needing to read the aws-infra repo.
</success_criteria>

<output>
After completion, create `.planning/quick/260327-vgp-add-production-deployment-section-to-rea/260327-vgp-SUMMARY.md`
</output>

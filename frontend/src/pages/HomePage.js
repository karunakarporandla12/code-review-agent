import React from "react";
import {
    Badge,
    Button,
    Card,
    Col,
    Container,
    Row,
} from "react-bootstrap";
import { Link } from "react-router-dom";

import HomeNavbar from "../components/HomeNavbar";
import "./HomePage.css";

const developmentStages = [
    {
        number: "01",
        icon: "bi-lightbulb",
        title: "Think",
        description:
            "Hanuma understands your idea, requirements, business goals, technical constraints, and expected user experience.",
        color: "blue",
    },
    {
        number: "02",
        icon: "bi-code-slash",
        title: "Build",
        description:
            "Generate structured, maintainable, and production-focused application code using your selected technology stack.",
        color: "green",
    },
    {
        number: "03",
        icon: "bi-check2-circle",
        title: "Test",
        description:
            "Validate generated code through automated tests, quality checks, edge-case analysis, and security scanning.",
        color: "orange",
    },
    {
        number: "04",
        icon: "bi-file-earmark-code",
        title: "Review",
        description:
            "Review architecture, source code, vulnerabilities, performance, maintainability, and implementation quality.",
        color: "purple",
    },
    {
        number: "05",
        icon: "bi-rocket-takeoff",
        title: "Deploy",
        description:
            "Prepare and deploy approved applications using containerized, cloud-ready, and automated delivery workflows.",
        color: "pink",
    },
    {
  
    number: "06",
    icon: "bi-activity",
    title: "Monitor",
    description:
        "Monitor the deployed application for runtime health, performance, errors, availability, logs, and production issues.",
    color: "cyan",
},
];

const capabilities = [
    {
        icon: "bi-chat-square-text",
        title: "Requirement Understanding",
        description:
            "Transform natural-language ideas and project requirements into structured development plans.",
        color: "blue",
    },
    {
        icon: "bi-braces",
        title: "Application Generation",
        description:
            "Generate frontend, backend, database, configuration, and deployment code from one centralized workspace.",
        color: "green",
    },
    {
        icon: "bi-bug",
        title: "Automated Testing",
        description:
            "Detect functional defects, validation gaps, edge cases, failed scenarios, and regression risks.",
        color: "orange",
    },
    {
        icon: "bi-shield-check",
        title: "Security Review",
        description:
            "Analyze authentication, authorization, exposed secrets, injection risks, and insecure implementation patterns.",
        color: "purple",
    },
    {
        icon: "bi-git",
        title: "Repository Intelligence",
        description:
            "Connect authorized repositories and analyze project structure, source files, dependencies, and changes.",
        color: "cyan",
    },
    {
        icon: "bi-cloud-check",
        title: "Deployment Readiness",
        description:
            "Prepare reviewed applications for Docker, Kubernetes, CI/CD, cloud, and production environments.",
        color: "pink",
    },
];

const securityFeatures = [
    {
        icon: "bi-lock",
        title: "End-to-End Encryption",
        description:
            "Protect source code and project information while data moves between supported application services.",
    },
    {
        icon: "bi-key",
        title: "Secure Key Handling",
        description:
            "Keep encryption keys and sensitive configuration outside normal source-code workflows.",
    },
    {
        icon: "bi-person-lock",
        title: "Controlled Access",
        description:
            "Restrict projects, repositories, reports, and deployment operations to authorized users.",
    },
    {
        icon: "bi-journal-check",
        title: "Activity Visibility",
        description:
            "Maintain visibility into analysis, reviews, generated output, and important project operations.",
    },
];

const codeLines = [
    {
        line: "01",
        content: "const project = await hanuma.understand(requirements);",
        type: "normal",
    },
    {
        line: "02",
        content: "const source = await hanuma.build(project);",
        type: "highlight",
    },
    {
        line: "03",
        content: "const tests = await hanuma.test(source);",
        type: "normal",
    },
    {
        line: "04",
        content: "const review = await hanuma.review(source);",
        type: "normal",
    },
    {
        line: "05",
        content: "await hanuma.deploy(review.approvedBuild);",
        type: "success",
    },
];

const SecurityItem = ({ icon, title, description }) => {
    return (
        <div className="hanuma-security-item h-100">
            <div className="hanuma-security-item-icon">
                <i className={`bi ${icon}`}></i>
            </div>

            <div>
                <h5>{title}</h5>
                <p>{description}</p>
            </div>
        </div>
    );
};

const HomePage = () => {
    return (
        <div className="hanuma-home-page">
            <HomeNavbar />

            <main>
                <section className="hanuma-hero">
                    <div className="hanuma-hero-grid"></div>
                    <div className="hanuma-orb hanuma-orb-one"></div>
                    <div className="hanuma-orb hanuma-orb-two"></div>
                    <div className="hanuma-orb hanuma-orb-three"></div>

                    <Container className="position-relative">
                        <Row className="align-items-center g-5">
                            <Col lg={6}>
                                <Badge className="hanuma-hero-badge">
                                    <span className="hanuma-badge-pulse"></span>
                                    Encrypted AI Development Agent
                                </Badge>

                                <div className="hanuma-brand-eyebrow">
                                    Super Intelligent Development
                                </div>

                                <h1 className="hanuma-hero-title">
                                    From idea to impact with
                                    <span> Hanuma</span>
                                </h1>

                                <p className="hanuma-hero-subtitle">
                                    Think. Build. Test. Review. Deploy.
                                </p>

                                <p className="hanuma-hero-description">
                                    Hanuma is an intelligent development partner
                                    that understands your idea, generates
                                    application code, validates implementation,
                                    reviews quality and security, and prepares
                                    your project for deployment.
                                </p>

                                <div className="hanuma-hero-actions">
                                    <Button
                                        as={Link}
                                        to="/register"
                                        size="lg"
                                        className="hanuma-primary-button"
                                    >
                                        Start Building
                                        <i className="bi bi-arrow-right ms-2"></i>
                                    </Button>

                                    <a
                                        href="#how-it-works"
                                        className="hanuma-secondary-button"
                                    >
                                        <i className="bi bi-play-circle me-2"></i>
                                        See How Hanuma Works
                                    </a>
                                </div>

                                <div className="hanuma-trust-row">
                                    <div>
                                        <i className="bi bi-lock-fill"></i>
                                        End-to-End Encrypted
                                    </div>

                                    <div>
                                        <i className="bi bi-shield-check"></i>
                                        Security Reviewed
                                    </div>

                                    <div>
                                        <i className="bi bi-rocket-takeoff"></i>
                                        Deployment Ready
                                    </div>
                                </div>
                            </Col>

                            <Col lg={6}>
                                <div className="hanuma-agent-visual">
                                    <div className="hanuma-agent-glow"></div>

                                    <div className="hanuma-agent-window">
                                        <div className="hanuma-window-header">
                                            <div className="hanuma-window-dots">
                                                <span></span>
                                                <span></span>
                                                <span></span>
                                            </div>

                                            <div className="hanuma-window-title">
                                                <i className="bi bi-stars"></i>
                                                Hanuma Agent
                                            </div>

                                            <Badge className="hanuma-live-badge">
                                                Active
                                            </Badge>
                                        </div>

                                        <div className="hanuma-agent-body">
                                            <div className="hanuma-agent-head">
                                                <div className="hanuma-agent-antenna">
                                                    <span></span>
                                                </div>

                                                <div className="hanuma-agent-face">
                                                    <div className="hanuma-eye"></div>
                                                    <div className="hanuma-eye"></div>
                                                </div>
                                            </div>

                                            <div className="hanuma-agent-name">
                                                HANUMA
                                            </div>

                                            <p>
                                                Your secure AI development partner
                                            </p>

                                            <div className="hanuma-code-panel">
                                                {codeLines.map((codeLine) => (
                                                    <div
                                                        className={`hanuma-code-line ${codeLine.type}`}
                                                        key={codeLine.line}
                                                    >
                                                        <span>
                                                            {codeLine.line}
                                                        </span>

                                                        <code>
                                                            {codeLine.content}
                                                        </code>
                                                    </div>
                                                ))}
                                            </div>

                                            <div className="hanuma-agent-status">
                                                <div>
                                                    <i className="bi bi-check-circle-fill"></i>
                                                    Code Generated
                                                </div>

                                                <div>
                                                    <i className="bi bi-check-circle-fill"></i>
                                                    Tests Passed
                                                </div>

                                                <div>
                                                    <i className="bi bi-check-circle-fill"></i>
                                                    Security Checked
                                                </div>

                                                <div>
                                                    <i className="bi bi-check-circle-fill"></i>
                                                    Ready to Deploy
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="hanuma-floating-card hanuma-floating-idea">
                                        <i className="bi bi-lightbulb"></i>

                                        <div>
                                            <span>From</span>
                                            <strong>Idea</strong>
                                        </div>
                                    </div>

                                    <div className="hanuma-floating-card hanuma-floating-deploy">
                                        <i className="bi bi-rocket-takeoff"></i>

                                        <div>
                                            <span>Ready for</span>
                                            <strong>Deployment</strong>
                                        </div>
                                    </div>

                                    <div className="hanuma-encryption-chip">
                                        <i className="bi bi-lock-fill"></i>

                                        <div>
                                            <strong>Encrypted</strong>
                                            <span>Secure processing</span>
                                        </div>
                                    </div>
                                </div>
                            </Col>
                        </Row>
                    </Container>
                </section>

                {/* <section
                    id="how-it-works"
                    className="hanuma-workflow-section"
                >
                    <Container>
                        <div className="hanuma-section-heading">
                            <span className="hanuma-section-label">
                                INTELLIGENCE IN ACTION
                            </span>

                            <h2>
                                One intelligent agent.
                                <br />
                                Your complete development lifecycle.
                            </h2>

                            <p>
                                Hanuma connects every important development
                                stage into one clear workflow, from
                                understanding the original idea to preparing
                                the application for deployment.
                            </p>
                        </div>

                        <div className="hanuma-workflow-container">
                            <div className="hanuma-workflow-line"></div>

                            <Row className="g-4">
                                {developmentStages.map((stage) => (
                                    <Col
                                        sm={12}
                                        md={6}
                                        lg
                                        key={stage.number}
                                    >
                                        <div className="hanuma-workflow-card h-100">
                                            <div className="hanuma-stage-number">
                                                {stage.number}
                                            </div>

                                            <div
                                                className={`hanuma-stage-icon ${stage.color}`}
                                            >
                                                <i
                                                    className={`bi ${stage.icon}`}
                                                ></i>
                                            </div>

                                            <h4>{stage.title}</h4>

                                            <p>
                                                {stage.description}
                                            </p>
                                        </div>
                                    </Col>
                                ))}
                            </Row>
                        </div>
                    </Container>
                </section> */}

                <section
    id="how-it-works"
    className="hanuma-workflow-section"
>
    <Container>
        <div className="hanuma-section-heading">
            <span className="hanuma-section-label">
                INTELLIGENCE IN ACTION
            </span>

            <h2>
                One intelligent agent.
                <br />
                Your complete development lifecycle.
            </h2>

            <p>
                Hanuma connects every important development stage
                into one clear workflow, from understanding the
                original idea to building, testing, reviewing,
                deploying, and completing the final review.
            </p>
        </div>

        <div className="hanuma-workflow-container">
            <div className="hanuma-workflow-line"></div>

            <Row className="g-4">
                {developmentStages.map((stage) => (
                    <Col
                        sm={12}
                        md={6}
                        lg={4}
                        key={stage.number}
                    >
                        <div className="hanuma-workflow-card h-100">
                            <div className="hanuma-stage-number">
                                {stage.number}
                            </div>

                            <div
                                className={`hanuma-stage-icon ${stage.color}`}
                            >
                                <i
                                    className={`bi ${stage.icon}`}
                                ></i>
                            </div>

                            <h4>
                                {stage.title}
                            </h4>

                            <p>
                                {stage.description}
                            </p>
                        </div>
                    </Col>
                ))}
            </Row>
        </div>
    </Container>
</section>

                <section
                    id="capabilities"
                    className="hanuma-capabilities-section"
                >
                    <Container>
                        <div className="hanuma-section-heading">
                            <span className="hanuma-section-label">
                                MORE THAN A CODE REVIEW TOOL
                            </span>

                            <h2>
                                Everything needed to build smarter software
                            </h2>

                            <p>
                                Move from requirements to working application
                                code while keeping testing, security, review,
                                quality, and deployment readiness connected.
                            </p>
                        </div>

                        <Row className="g-4">
                            {capabilities.map((capability) => (
                                <Col
                                    sm={12}
                                    md={6}
                                    lg={4}
                                    key={capability.title}
                                >
                                    <Card className="hanuma-capability-card h-100">
                                        <Card.Body>
                                            <div
                                                className={`hanuma-capability-icon ${capability.color}`}
                                            >
                                                <i
                                                    className={`bi ${capability.icon}`}
                                                ></i>
                                            </div>

                                            <Card.Title>
                                                {capability.title}
                                            </Card.Title>

                                            <Card.Text>
                                                {capability.description}
                                            </Card.Text>

                                            <div className="hanuma-card-link">
                                                Explore capability
                                                <i className="bi bi-arrow-right"></i>
                                            </div>
                                        </Card.Body>
                                    </Card>
                                </Col>
                            ))}
                        </Row>
                    </Container>
                </section>

                <section
                    id="security"
                    className="hanuma-encryption-section"
                >
                    <Container>
                        <Row className="align-items-center g-5">
                            <Col lg={6}>
                                <div className="hanuma-encryption-visual">
                                    <div className="hanuma-encryption-ring ring-one"></div>

                                    <div className="hanuma-encryption-ring ring-two"></div>

                                    <div className="hanuma-encryption-ring ring-three"></div>

                                    <div className="hanuma-encryption-core">
                                        <i className="bi bi-shield-lock-fill"></i>

                                        <strong>
                                            Encrypted
                                        </strong>

                                        <span>
                                            Protected workflow
                                        </span>
                                    </div>

                                    <div className="hanuma-secure-node node-code">
                                        <i className="bi bi-code-square"></i>
                                        Source Code
                                    </div>

                                    <div className="hanuma-secure-node node-repository">
                                        <i className="bi bi-git"></i>
                                        Repository
                                    </div>

                                    <div className="hanuma-secure-node node-report">
                                        <i className="bi bi-file-earmark-lock"></i>
                                        Review Report
                                    </div>

                                    <div className="hanuma-secure-node node-deployment">
                                        <i className="bi bi-cloud-lock"></i>
                                        Deployment
                                    </div>
                                </div>
                            </Col>

                            <Col lg={6}>
                                <span className="hanuma-section-label">
                                    SECURITY ACROSS THE WORKFLOW
                                </span>

                                <h2 className="hanuma-encryption-title">
                                    Your code remains protected throughout
                                    the development lifecycle
                                </h2>

                                <p className="hanuma-encryption-description">
                                    Hanuma places encryption, access control,
                                    secure key management, and activity
                                    visibility around important development
                                    operations.
                                </p>

                                <Row className="g-3 mt-4">
                                    {securityFeatures.map((feature) => (
                                        <Col
                                            md={6}
                                            key={feature.title}
                                        >
                                            <SecurityItem
                                                icon={feature.icon}
                                                title={feature.title}
                                                description={feature.description}
                                            />
                                        </Col>
                                    ))}
                                </Row>

                                <div className="hanuma-security-summary">
                                    <div className="hanuma-security-summary-icon">
                                        <i className="bi bi-shield-check"></i>
                                    </div>

                                    <div>
                                        <strong>
                                            Secure by design
                                        </strong>

                                        <p>
                                            Security controls are integrated
                                            into code generation, analysis,
                                            review, and deployment preparation.
                                        </p>
                                    </div>
                                </div>
                            </Col>
                        </Row>
                    </Container>
                </section>

                <section className="hanuma-cta-section">
                    <Container>
                        <div className="hanuma-cta-card">
                            <div className="hanuma-cta-grid"></div>

                            <div className="hanuma-cta-orb hanuma-cta-orb-one"></div>

                            <div className="hanuma-cta-orb hanuma-cta-orb-two"></div>

                            <div className="hanuma-cta-content">
                                <Badge className="hanuma-cta-badge">
                                    <i className="bi bi-stars me-2"></i>
                                    BUILD WITH HANUMA
                                </Badge>

                                <h2>
                                    Turn your next idea into secure,
                                    production-focused software
                                </h2>

                                <p>
                                    Describe the application you want to build
                                    and let Hanuma connect planning, code
                                    generation, testing, review, security, and
                                    deployment preparation.
                                </p>

                                <div className="hanuma-cta-actions">
                                    <Button
                                        as={Link}
                                        to="/register"
                                        size="lg"
                                        variant="light"
                                        className="hanuma-cta-primary-button"
                                    >
                                        Start Building
                                        <i className="bi bi-arrow-right ms-2"></i>
                                    </Button>

                                    <Button
                                        as={Link}
                                        to="/login"
                                        size="lg"
                                        variant="outline-light"
                                        className="hanuma-cta-secondary-button"
                                    >
                                        Sign In
                                    </Button>
                                </div>

                                <div className="hanuma-cta-trust">
                                    <span>
                                        <i className="bi bi-lock-fill"></i>
                                        Encrypted Workflow
                                    </span>

                                    <span>
                                        <i className="bi bi-check-circle-fill"></i>
                                        Reviewed Output
                                    </span>

                                    <span>
                                        <i className="bi bi-cloud-check-fill"></i>
                                        Deployment Ready
                                    </span>
                                </div>
                            </div>
                        </div>
                    </Container>
                </section>
            </main>

            <footer className="hanuma-footer">
                <Container>
                    <Row className="align-items-center gy-4">
                        <Col lg={5}>
                            <Link
                                to="/"
                                className="hanuma-footer-brand"
                            >
                                <div className="hanuma-footer-brand-icon">
                                    <i className="bi bi-stars"></i>
                                </div>

                                <div>
                                    <div className="hanuma-footer-brand-title">
                                        Hanuma
                                    </div>

                                    <div className="hanuma-footer-brand-subtitle">
                                        Super Intelligent Development Agent
                                    </div>
                                </div>
                            </Link>
                        </Col>

                        <Col lg={7}>
                            <div className="hanuma-footer-links">
                                <a href="/#how-it-works">
                                    How It Works
                                </a>

                                <a href="/#capabilities">
                                    Capabilities
                                </a>

                                <a href="/#security">
                                    Security
                                </a>

                                <Link to="/login">
                                    Sign In
                                </Link>

                                <Link to="/register">
                                    Get Started
                                </Link>
                            </div>
                        </Col>
                    </Row>

                    <div className="hanuma-footer-divider"></div>

                    <div className="hanuma-footer-bottom">
                        <span>
                            © {new Date().getFullYear()} Hanuma. All rights reserved.
                        </span>

                        <div className="hanuma-footer-status">
                            <span className="hanuma-status-dot"></span>
                            System Operational
                        </div>
                    </div>
                </Container>
            </footer>
        </div>
    );
};

export default HomePage;


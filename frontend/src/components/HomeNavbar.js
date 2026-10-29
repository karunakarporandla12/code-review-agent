// import React, { useEffect, useState } from "react";
// import { Button, Container, Navbar } from "react-bootstrap";
// import { Link, useLocation } from "react-router-dom";
// import "./HomeNavbar.css";

// const HomeNavbar = () => {
//     const location = useLocation();

//     const [scrolled, setScrolled] = useState(false);
//     const [expanded, setExpanded] = useState(false);

//     useEffect(() => {
//         const handleScroll = () => {
//             setScrolled(window.scrollY > 20);
//         };

//         handleScroll();

//         window.addEventListener("scroll", handleScroll);

//         return () => {
//             window.removeEventListener("scroll", handleScroll);
//         };
//     }, []);

//     const handleSectionClick = (event, sectionId) => {
//         setExpanded(false);

//         if (location.pathname === "/") {
//             event.preventDefault();

//             const section = document.getElementById(sectionId);

//             if (section) {
//                 section.scrollIntoView({
//                     behavior: "smooth",
//                     block: "start",
//                 });
//             }
//         }
//     };

//     const closeNavbar = () => {
//         setExpanded(false);
//     };

//     return (
//         <Navbar
//             expand="lg"
//             expanded={expanded}
//             onToggle={(isExpanded) => setExpanded(isExpanded)}
//             className={`home-navbar ${
//                 scrolled ? "home-navbar-scrolled" : ""
//             }`}
//         >
//             <Container>
//                 {/* <Navbar.Brand
//                     as={Link}
//                     to="/"
//                     onClick={handleBrandClick}
//                     className="home-navbar-brand"
//                 >
//                     <div className="navbar-logo">
//                         <i className="bi bi-code-slash"></i>
//                     </div>

//                     <div className="navbar-brand-content">
//                         <div className="navbar-brand-title">
//                             CodeReview
//                         </div>

//                         <div className="navbar-brand-subtitle">
//                             Intelligent Code Analysis
//                         </div>
//                     </div>
//                 </Navbar.Brand> */}

//                 <Navbar.Brand
//                     as={Link}
//                     to="/"
//                     className="d-flex align-items-center"
//                 >
//                     <div className="brand-icon me-2">
//                         <i className="bi bi-code-slash" />
//                     </div>

//                     <div>
//                         <div className="brand-title">
//                             Hanuma
//                         </div>

//                         <div className="brand-subtitle">
//                             Super Intelligent Code Analysis
//                         </div>
//                     </div>
//                 </Navbar.Brand>

//                 <Navbar.Toggle
//                     aria-controls="home-navbar-navigation"
//                     className="home-navbar-toggler"
//                 >
//                     <i className="bi bi-list"></i>
//                 </Navbar.Toggle>

//                 <Navbar.Collapse id="home-navbar-navigation">
//                     <div className="navbar-center-menu mx-auto">
//                         <Link
//                             to="/#features"
//                             className="home-nav-link"
//                             onClick={(event) =>
//                                 handleSectionClick(event, "features")
//                             }
//                         >
//                             <i className="bi bi-grid"></i>
//                             <span>Features</span>
//                         </Link>

//                         <Link
//                             to="/#workflow"
//                             className="home-nav-link"
//                             onClick={(event) =>
//                                 handleSectionClick(event, "workflow")
//                             }
//                         >
//                             <i className="bi bi-diagram-3"></i>
//                             <span>How It Works</span>
//                         </Link>

//                         <Link
//                             to="/#security"
//                             className="home-nav-link"
//                             onClick={(event) =>
//                                 handleSectionClick(event, "security")
//                             }
//                         >
//                             <i className="bi bi-shield-check"></i>
//                             <span>Security</span>
//                         </Link>
//                     </div>

//                     <div className="navbar-actions">
//                         <Button
//                             as={Link}
//                             to="/login"
//                             variant="link"
//                             className="navbar-signin-button"
//                             onClick={closeNavbar}
//                         >
//                             <i className="bi bi-box-arrow-in-right"></i>
//                             <span>Sign In</span>
//                         </Button>

//                         <Button
//                             as={Link}
//                             to="/register"
//                             className="navbar-get-started-button"
//                             onClick={closeNavbar}
//                         >
//                             <span>Get Started</span>
//                             <i className="bi bi-arrow-right"></i>
//                         </Button>
//                     </div>
//                 </Navbar.Collapse>
//             </Container>
//         </Navbar>
//     );
// };

// export default HomeNavbar;

import React, { useEffect, useState } from "react";
import {
    Button,
    Container,
    Nav,
    Navbar,
} from "react-bootstrap";
import {
    Link,
    useLocation,
    useNavigate,
} from "react-router-dom";
import "./HomeNavbar.css";

const HomeNavbar = () => {
    const location = useLocation();
    const navigate = useNavigate();

    const [scrolled, setScrolled] = useState(false);
    const [expanded, setExpanded] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            setScrolled(window.scrollY > 20);
        };

        handleScroll();

        window.addEventListener("scroll", handleScroll);

        return () => {
            window.removeEventListener("scroll", handleScroll);
        };
    }, []);

    const scrollToSection = (sectionId) => {
        window.setTimeout(() => {
            const section = document.getElementById(sectionId);

            if (section) {
                section.scrollIntoView({
                    behavior: "smooth",
                    block: "start",
                });
            }
        }, 100);
    };

    const handleSectionClick = (event, sectionId) => {
        event.preventDefault();

        setExpanded(false);

        if (location.pathname === "/") {
            scrollToSection(sectionId);
            return;
        }

        navigate(`/#${sectionId}`);

        window.setTimeout(() => {
            scrollToSection(sectionId);
        }, 150);
    };

    const handleBrandClick = () => {
        setExpanded(false);

        if (location.pathname === "/") {
            window.scrollTo({
                top: 0,
                behavior: "smooth",
            });
        }
    };

    const closeNavbar = () => {
        setExpanded(false);
    };

    return (
        <Navbar
            expand="lg"
            expanded={expanded}
            onToggle={(isExpanded) => setExpanded(isExpanded)}
            className={`home-navbar ${
                scrolled ? "home-navbar-scrolled" : ""
            }`}
        >
            <Container>
                <Navbar.Brand
                    as={Link}
                    to="/"
                    onClick={handleBrandClick}
                    className="home-navbar-brand"
                >
                    <div className="home-navbar-logo">
                        <i className="bi bi-stars"></i>
                    </div>

                    <div className="home-navbar-brand-content">
                        <div className="home-navbar-brand-title">
                            Hanuma
                        </div>

                        <div className="home-navbar-brand-subtitle">
                            Super Intelligent Development Agent
                        </div>
                    </div>
                </Navbar.Brand>

                <Navbar.Toggle
                    aria-controls="home-navbar-navigation"
                    className="home-navbar-toggler"
                >
                    <i className="bi bi-list"></i>
                </Navbar.Toggle>

                <Navbar.Collapse id="home-navbar-navigation">
                    <Nav className="navbar-center-menu mx-auto">
                        <a
                            href="/#capabilities"
                            className="home-nav-link"
                            onClick={(event) =>
                                handleSectionClick(
                                    event,
                                    "capabilities"
                                )
                            }
                        >
                            <i className="bi bi-grid"></i>
                            <span>Capabilities</span>
                        </a>

                        <a
                            href="/#how-it-works"
                            className="home-nav-link"
                            onClick={(event) =>
                                handleSectionClick(
                                    event,
                                    "how-it-works"
                                )
                            }
                        >
                            <i className="bi bi-diagram-3"></i>
                            <span>How It Works</span>
                        </a>

                        <a
                            href="/#security"
                            className="home-nav-link"
                            onClick={(event) =>
                                handleSectionClick(
                                    event,
                                    "security"
                                )
                            }
                        >
                            <i className="bi bi-shield-check"></i>
                            <span>Security</span>
                        </a>
                    </Nav>

                    <div className="navbar-actions">
                        <Button
                            as={Link}
                            to="/login"
                            variant="link"
                            className="navbar-signin-button"
                            onClick={closeNavbar}
                        >
                            <i className="bi bi-box-arrow-in-right"></i>
                            <span>Sign In</span>
                        </Button>

                        <Button
                            as={Link}
                            to="/register"
                            className="navbar-get-started-button"
                            onClick={closeNavbar}
                        >
                            <span>Get Started</span>
                            <i className="bi bi-arrow-right"></i>
                        </Button>
                    </div>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
};

export default HomeNavbar;
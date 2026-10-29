import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import api from "../api";
import { useAuth } from "../auth/AuthContext";
import { User } from "../types";
const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { setSession } = useAuth();
  const nav = useNavigate();
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      const r = await api.post<{ token: string; user: User }>(
        "/api/auth/login",
        { email, password },
      );
      setSession(r.data.token, r.data.user);
      nav("/dashboard");
    } catch (err: any) {
      setError(err.response?.data?.error || "Login failed");
    }
  };
  return (
    <div className="auth-card card shadow-sm p-4">
      <h2>Sign in</h2>
      <p className="text-muted">
        Review uploaded projects and Git repositories.
      </p>
      {error && <div className="alert alert-danger">{error}</div>}
      <form onSubmit={submit}>
        <label className="form-label">Email</label>
        <input
          className="form-control mb-3"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <label className="form-label">Password</label>
        <input
          className="form-control mb-3"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button className="btn btn-primary w-100">Sign in</button>
      </form>
      <p className="mt-3 mb-0">
        New user? <Link to="/register">Create account</Link>
      </p>
    </div>
  );
};
export default LoginPage;

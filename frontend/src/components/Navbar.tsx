import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Navbar.css';

const Navbar: React.FC = () => {
  const location = useLocation();

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-brand">
          <span className="brand-icon">🎯</span>
          CodeRoulette
        </Link>
        
        <div className="navbar-menu">
          <Link 
            to="/" 
            className={`navbar-link ${location.pathname === '/' ? 'active' : ''}`}
          >
            Home
          </Link>
          <Link 
            to="/battle" 
            className={`navbar-link ${location.pathname === '/battle' ? 'active' : ''}`}
          >
            Battle
          </Link>
          <Link 
            to="/leaderboard" 
            className={`navbar-link ${location.pathname === '/leaderboard' ? 'active' : ''}`}
          >
            Leaderboard
          </Link>
          <Link 
            to="/profile" 
            className={`navbar-link ${location.pathname === '/profile' ? 'active' : ''}`}
          >
            Profile
          </Link>
        </div>
        
        <div className="navbar-user">
          <div className="user-info">
            <span className="user-avatar">👤</span>
            <span className="user-name">Player</span>
            <span className="user-rating">1200</span>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

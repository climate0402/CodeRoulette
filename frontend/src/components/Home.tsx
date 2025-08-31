import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const [difficulty, setDifficulty] = useState('medium');
  const [language, setLanguage] = useState('go');
  const [isSearching, setIsSearching] = useState(false);

  const handleStartBattle = async () => {
    setIsSearching(true);
    
    // Simulate matchmaking
    setTimeout(() => {
      setIsSearching(false);
      navigate('/battle');
    }, 2000);
  };

  return (
    <div className="home">
      <div className="hero-section">
        <div className="hero-content">
          <h1 className="hero-title">
            Welcome to <span className="highlight">CodeRoulette</span>
          </h1>
          <p className="hero-subtitle">
            Real-time programming battles where code meets competition
          </p>
          
          <div className="battle-setup">
            <div className="setup-options">
              <div className="option-group">
                <label>Difficulty</label>
                <select 
                  value={difficulty} 
                  onChange={(e) => setDifficulty(e.target.value)}
                  className="option-select"
                >
                  <option value="easy">Easy</option>
                  <option value="medium">Medium</option>
                  <option value="hard">Hard</option>
                </select>
              </div>
              
              <div className="option-group">
                <label>Language</label>
                <select 
                  value={language} 
                  onChange={(e) => setLanguage(e.target.value)}
                  className="option-select"
                >
                  <option value="go">Go</option>
                  <option value="python">Python</option>
                  <option value="javascript">JavaScript</option>
                </select>
              </div>
            </div>
            
            <button 
              className="start-battle-btn"
              onClick={handleStartBattle}
              disabled={isSearching}
            >
              {isSearching ? (
                <>
                  <span className="spinner"></span>
                  Finding Opponent...
                </>
              ) : (
                'Start Battle'
              )}
            </button>
          </div>
        </div>
        
        <div className="hero-visual">
          <div className="code-preview">
            <div className="code-header">
              <div className="code-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span className="code-title">battle.go</span>
            </div>
            <div className="code-content">
              <div className="code-line">
                <span className="keyword">func</span> <span className="function">solve</span><span className="bracket">(</span><span className="param">input</span> <span className="type">string</span><span className="bracket">)</span> <span className="type">string</span> <span className="bracket">{'{'}</span>
              </div>
              <div className="code-line">
                &nbsp;&nbsp;&nbsp;&nbsp;<span className="comment">// Your code here</span>
              </div>
              <div className="code-line">
                &nbsp;&nbsp;&nbsp;&nbsp;<span className="keyword">return</span> <span className="string">"solution"</span>
              </div>
              <div className="code-line">
                <span className="bracket">{'}'}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div className="features-section">
        <h2 className="section-title">Features</h2>
        <div className="features-grid">
          <div className="feature-card">
            <div className="feature-icon">⚡</div>
            <h3>Real-time Battles</h3>
            <p>Compete against other programmers in live coding challenges</p>
          </div>
          
          <div className="feature-card">
            <div className="feature-icon">🎯</div>
            <h3>Skill Cards</h3>
            <p>Use special abilities to gain advantages during battles</p>
          </div>
          
          <div className="feature-card">
            <div className="feature-icon">📊</div>
            <h3>Detailed Reports</h3>
            <p>Get comprehensive battle reports and performance analytics</p>
          </div>
          
          <div className="feature-card">
            <div className="feature-icon">👥</div>
            <h3>Spectator Mode</h3>
            <p>Watch other battles and learn from top players</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;

import React, { useState, useEffect } from 'react';
import './Profile.css';

interface UserStats {
  id: string;
  username: string;
  email: string;
  rating: number;
  wins: number;
  losses: number;
  totalMatches: number;
  winRate: number;
  averageScore: number;
  bestScore: number;
  joinDate: string;
  favoriteLanguage: string;
  skillCards: number;
}

interface MatchHistory {
  id: string;
  opponent: string;
  problem: string;
  result: 'win' | 'loss';
  score: number;
  duration: number;
  date: string;
}

const Profile: React.FC = () => {
  const [userStats, setUserStats] = useState<UserStats | null>(null);
  const [matchHistory, setMatchHistory] = useState<MatchHistory[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'stats' | 'history' | 'achievements'>('stats');

  useEffect(() => {
    // Simulate API call
    setTimeout(() => {
      const sampleStats: UserStats = {
        id: '1',
        username: 'CodeMaster',
        email: 'codemaster@example.com',
        rating: 1850,
        wins: 45,
        losses: 12,
        totalMatches: 57,
        winRate: 78.9,
        averageScore: 87.5,
        bestScore: 100,
        joinDate: '2024-01-15',
        favoriteLanguage: 'Go',
        skillCards: 12
      };

      const sampleHistory: MatchHistory[] = [
        {
          id: '1',
          opponent: 'AlgorithmKing',
          problem: 'Two Sum',
          result: 'win',
          score: 95,
          duration: 180,
          date: '2024-01-20'
        },
        {
          id: '2',
          opponent: 'DataNinja',
          problem: 'Valid Parentheses',
          result: 'loss',
          score: 60,
          duration: 240,
          date: '2024-01-19'
        },
        {
          id: '3',
          opponent: 'BugHunter',
          problem: 'Reverse String',
          result: 'win',
          score: 100,
          duration: 120,
          date: '2024-01-18'
        },
        {
          id: '4',
          opponent: 'SpeedCoder',
          problem: 'Merge Sorted Arrays',
          result: 'win',
          score: 88,
          duration: 200,
          date: '2024-01-17'
        },
        {
          id: '5',
          opponent: 'LogicWizard',
          problem: 'Binary Search',
          result: 'loss',
          score: 45,
          duration: 300,
          date: '2024-01-16'
        }
      ];

      setUserStats(sampleStats);
      setMatchHistory(sampleHistory);
      setLoading(false);
    }, 1000);
  }, []);

  const getRatingColor = (rating: number) => {
    if (rating >= 1800) return '#ff6b6b';
    if (rating >= 1600) return '#4ecdc4';
    if (rating >= 1400) return '#45b7d1';
    return '#96ceb4';
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (loading) {
    return (
      <div className="profile-loading">
        <div className="loading-spinner"></div>
        <p>Loading profile...</p>
      </div>
    );
  }

  if (!userStats) {
    return (
      <div className="profile-error">
        <p>Failed to load profile data</p>
      </div>
    );
  }

  return (
    <div className="profile">
      <div className="profile-header">
        <div className="profile-avatar">
          {userStats.username[0].toUpperCase()}
        </div>
        <div className="profile-info">
          <h1 className="profile-name">{userStats.username}</h1>
          <div className="profile-rating">
            <span 
              className="rating-value"
              style={{ color: getRatingColor(userStats.rating) }}
            >
              {userStats.rating}
            </span>
            <span className="rating-label">Rating</span>
          </div>
          <p className="profile-email">{userStats.email}</p>
        </div>
        <div className="profile-stats-summary">
          <div className="stat-item">
            <div className="stat-value">{userStats.wins}</div>
            <div className="stat-label">Wins</div>
          </div>
          <div className="stat-item">
            <div className="stat-value">{userStats.losses}</div>
            <div className="stat-label">Losses</div>
          </div>
          <div className="stat-item">
            <div className="stat-value">{userStats.winRate}%</div>
            <div className="stat-label">Win Rate</div>
          </div>
        </div>
      </div>

      <div className="profile-tabs">
        <button
          className={`tab-btn ${activeTab === 'stats' ? 'active' : ''}`}
          onClick={() => setActiveTab('stats')}
        >
          Statistics
        </button>
        <button
          className={`tab-btn ${activeTab === 'history' ? 'active' : ''}`}
          onClick={() => setActiveTab('history')}
        >
          Match History
        </button>
        <button
          className={`tab-btn ${activeTab === 'achievements' ? 'active' : ''}`}
          onClick={() => setActiveTab('achievements')}
        >
          Achievements
        </button>
      </div>

      <div className="profile-content">
        {activeTab === 'stats' && (
          <div className="stats-content">
            <div className="stats-grid">
              <div className="stat-card">
                <h3>Total Matches</h3>
                <div className="stat-number">{userStats.totalMatches}</div>
              </div>
              <div className="stat-card">
                <h3>Average Score</h3>
                <div className="stat-number">{userStats.averageScore}</div>
              </div>
              <div className="stat-card">
                <h3>Best Score</h3>
                <div className="stat-number">{userStats.bestScore}</div>
              </div>
              <div className="stat-card">
                <h3>Favorite Language</h3>
                <div className="stat-number">{userStats.favoriteLanguage}</div>
              </div>
              <div className="stat-card">
                <h3>Skill Cards</h3>
                <div className="stat-number">{userStats.skillCards}</div>
              </div>
              <div className="stat-card">
                <h3>Member Since</h3>
                <div className="stat-number">{formatDate(userStats.joinDate)}</div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'history' && (
          <div className="history-content">
            <div className="match-history">
              {matchHistory.map((match) => (
                <div key={match.id} className="match-item">
                  <div className="match-result">
                    <span className={`result-badge ${match.result}`}>
                      {match.result.toUpperCase()}
                    </span>
                  </div>
                  <div className="match-details">
                    <div className="match-opponent">vs {match.opponent}</div>
                    <div className="match-problem">{match.problem}</div>
                    <div className="match-meta">
                      Score: {match.score} • Duration: {formatDuration(match.duration)} • {formatDate(match.date)}
                    </div>
                  </div>
                  <div className="match-score">
                    <div className="score-value">{match.score}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'achievements' && (
          <div className="achievements-content">
            <div className="achievements-grid">
              <div className="achievement-card earned">
                <div className="achievement-icon">🏆</div>
                <h3>First Victory</h3>
                <p>Win your first battle</p>
                <div className="achievement-date">Earned Jan 15, 2024</div>
              </div>
              <div className="achievement-card earned">
                <div className="achievement-icon">⚡</div>
                <h3>Speed Demon</h3>
                <p>Complete a battle in under 2 minutes</p>
                <div className="achievement-date">Earned Jan 18, 2024</div>
              </div>
              <div className="achievement-card earned">
                <div className="achievement-icon">🎯</div>
                <h3>Perfect Score</h3>
                <p>Get 100% on a problem</p>
                <div className="achievement-date">Earned Jan 18, 2024</div>
              </div>
              <div className="achievement-card">
                <div className="achievement-icon">🔥</div>
                <h3>Win Streak</h3>
                <p>Win 10 battles in a row</p>
                <div className="achievement-progress">7/10</div>
              </div>
              <div className="achievement-card">
                <div className="achievement-icon">💎</div>
                <h3>Diamond Rank</h3>
                <p>Reach 2000 rating</p>
                <div className="achievement-progress">1850/2000</div>
              </div>
              <div className="achievement-card">
                <div className="achievement-icon">🎮</div>
                <h3>Skill Master</h3>
                <p>Use 50 skill cards</p>
                <div className="achievement-progress">12/50</div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Profile;

import React, { useState, useEffect } from 'react';
import './Leaderboard.css';

interface Player {
  id: string;
  username: string;
  rating: number;
  wins: number;
  losses: number;
  winRate: number;
}

const Leaderboard: React.FC = () => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [timeFilter, setTimeFilter] = useState<'all' | 'week' | 'month'>('all');

  useEffect(() => {
    // Simulate API call
    setTimeout(() => {
      const samplePlayers: Player[] = [
        {
          id: '1',
          username: 'CodeMaster',
          rating: 1850,
          wins: 45,
          losses: 12,
          winRate: 78.9
        },
        {
          id: '2',
          username: 'AlgorithmKing',
          rating: 1820,
          wins: 38,
          losses: 15,
          winRate: 71.7
        },
        {
          id: '3',
          username: 'DataNinja',
          rating: 1790,
          wins: 42,
          losses: 18,
          winRate: 70.0
        },
        {
          id: '4',
          username: 'BugHunter',
          rating: 1760,
          wins: 35,
          losses: 20,
          winRate: 63.6
        },
        {
          id: '5',
          username: 'SpeedCoder',
          rating: 1730,
          wins: 40,
          losses: 25,
          winRate: 61.5
        },
        {
          id: '6',
          username: 'LogicWizard',
          rating: 1700,
          wins: 33,
          losses: 22,
          winRate: 60.0
        },
        {
          id: '7',
          username: 'ByteBuster',
          rating: 1670,
          wins: 28,
          losses: 18,
          winRate: 60.9
        },
        {
          id: '8',
          username: 'StackOverflow',
          rating: 1640,
          wins: 31,
          losses: 24,
          winRate: 56.4
        },
        {
          id: '9',
          username: 'FunctionCall',
          rating: 1610,
          wins: 26,
          losses: 20,
          winRate: 56.5
        },
        {
          id: '10',
          username: 'LoopMaster',
          rating: 1580,
          wins: 29,
          losses: 26,
          winRate: 52.7
        }
      ];
      setPlayers(samplePlayers);
      setLoading(false);
    }, 1000);
  }, [timeFilter]);

  const getRankIcon = (index: number) => {
    switch (index) {
      case 0:
        return '🥇';
      case 1:
        return '🥈';
      case 2:
        return '🥉';
      default:
        return `#${index + 1}`;
    }
  };

  const getRatingColor = (rating: number) => {
    if (rating >= 1800) return '#ff6b6b';
    if (rating >= 1600) return '#4ecdc4';
    if (rating >= 1400) return '#45b7d1';
    return '#96ceb4';
  };

  if (loading) {
    return (
      <div className="leaderboard-loading">
        <div className="loading-spinner"></div>
        <p>Loading leaderboard...</p>
      </div>
    );
  }

  return (
    <div className="leaderboard">
      <div className="leaderboard-header">
        <h1>Leaderboard</h1>
        <div className="time-filters">
          <button
            className={`filter-btn ${timeFilter === 'all' ? 'active' : ''}`}
            onClick={() => setTimeFilter('all')}
          >
            All Time
          </button>
          <button
            className={`filter-btn ${timeFilter === 'week' ? 'active' : ''}`}
            onClick={() => setTimeFilter('week')}
          >
            This Week
          </button>
          <button
            className={`filter-btn ${timeFilter === 'month' ? 'active' : ''}`}
            onClick={() => setTimeFilter('month')}
          >
            This Month
          </button>
        </div>
      </div>

      <div className="leaderboard-content">
        <div className="leaderboard-table">
          <div className="table-header">
            <div className="rank-col">Rank</div>
            <div className="player-col">Player</div>
            <div className="rating-col">Rating</div>
            <div className="wins-col">Wins</div>
            <div className="losses-col">Losses</div>
            <div className="rate-col">Win Rate</div>
          </div>

          {players.map((player, index) => (
            <div key={player.id} className="table-row">
              <div className="rank-col">
                <span className="rank-icon">{getRankIcon(index)}</span>
              </div>
              <div className="player-col">
                <div className="player-info">
                  <div className="player-avatar">
                    {player.username[0].toUpperCase()}
                  </div>
                  <span className="player-name">{player.username}</span>
                </div>
              </div>
              <div className="rating-col">
                <span 
                  className="rating-value"
                  style={{ color: getRatingColor(player.rating) }}
                >
                  {player.rating}
                </span>
              </div>
              <div className="wins-col">
                <span className="wins-value">{player.wins}</span>
              </div>
              <div className="losses-col">
                <span className="losses-value">{player.losses}</span>
              </div>
              <div className="rate-col">
                <div className="win-rate">
                  <span className="rate-value">{player.winRate}%</span>
                  <div className="rate-bar">
                    <div 
                      className="rate-fill"
                      style={{ width: `${player.winRate}%` }}
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="leaderboard-stats">
          <div className="stats-card">
            <h3>Total Players</h3>
            <div className="stat-value">1,247</div>
          </div>
          <div className="stats-card">
            <h3>Active Today</h3>
            <div className="stat-value">89</div>
          </div>
          <div className="stats-card">
            <h3>Battles Today</h3>
            <div className="stat-value">156</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Leaderboard;

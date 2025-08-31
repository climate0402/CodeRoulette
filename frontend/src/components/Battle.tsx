import React, { useState, useEffect, useRef } from 'react';
import Editor from '@monaco-editor/react';
import { io, Socket } from 'socket.io-client';
import './Battle.css';

interface Problem {
  id: string;
  title: string;
  description: string;
  difficulty: string;
  language: string;
  testCases: TestCase[];
}

interface TestCase {
  input: string;
  output: string;
}

interface Player {
  id: string;
  username: string;
  status: 'online' | 'offline' | 'coding';
  score: number;
  submissions: number;
}

interface SkillCard {
  id: string;
  name: string;
  description: string;
  cost: number;
  rarity: string;
}

const Battle: React.FC = () => {
  const [problem, setProblem] = useState<Problem | null>(null);
  const [player1, setPlayer1] = useState<Player>({
    id: '1',
    username: 'Player 1',
    status: 'online',
    score: 0,
    submissions: 0
  });
  const [player2, setPlayer2] = useState<Player>({
    id: '2',
    username: 'Player 2',
    status: 'online',
    score: 0,
    submissions: 0
  });
  const [code, setCode] = useState('');
  const [timeLeft, setTimeLeft] = useState(300); // 5 minutes
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [skillCards, setSkillCards] = useState<SkillCard[]>([]);
  const [socket, setSocket] = useState<Socket | null>(null);
  const [matchStatus, setMatchStatus] = useState<'waiting' | 'active' | 'completed'>('active');

  const editorRef = useRef<any>(null);

  useEffect(() => {
    // Initialize WebSocket connection
    const newSocket = io('ws://localhost:8080');
    setSocket(newSocket);

    // Load sample problem
    const sampleProblem: Problem = {
      id: '1',
      title: 'Two Sum',
      description: 'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.\n\nYou can return the answer in any order.',
      difficulty: 'easy',
      language: 'go',
      testCases: [
        { input: '[2,7,11,15], 9', output: '[0,1]' },
        { input: '[3,2,4], 6', output: '[1,2]' },
        { input: '[3,3], 6', output: '[0,1]' }
      ]
    };
    setProblem(sampleProblem);

    // Load sample skill cards
    const sampleSkillCards: SkillCard[] = [
      {
        id: '1',
        name: 'Code Peek',
        description: 'View one line of opponent\'s code',
        cost: 1,
        rarity: 'common'
      },
      {
        id: '2',
        name: 'Hint',
        description: 'Get a hint for the problem',
        cost: 1,
        rarity: 'common'
      },
      {
        id: '3',
        name: 'Time Boost',
        description: 'Get 30 seconds extra time',
        cost: 2,
        rarity: 'rare'
      }
    ];
    setSkillCards(sampleSkillCards);

    // Set initial code
    setCode(`func twoSum(nums []int, target int) []int {
    // Your code here
    return []int{}
}`);

    return () => {
      newSocket.close();
    };
  }, []);

  useEffect(() => {
    // Timer countdown
    if (timeLeft > 0 && matchStatus === 'active') {
      const timer = setTimeout(() => {
        setTimeLeft(timeLeft - 1);
      }, 1000);
      return () => clearTimeout(timer);
    } else if (timeLeft === 0) {
      setMatchStatus('completed');
    }
  }, [timeLeft, matchStatus]);

  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor;
  };

  const handleCodeChange = (value: string | undefined) => {
    if (value !== undefined) {
      setCode(value);
    }
  };

  const handleSubmit = async () => {
    if (!code.trim() || isSubmitting) return;

    setIsSubmitting(true);
    
    // Simulate submission
    setTimeout(() => {
      setIsSubmitting(false);
      setPlayer1(prev => ({
        ...prev,
        submissions: prev.submissions + 1,
        score: Math.min(100, prev.score + 25)
      }));
    }, 2000);
  };

  const handleSkillCardUse = (card: SkillCard) => {
    // Simulate skill card usage
    console.log(`Using skill card: ${card.name}`);
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (!problem) {
    return (
      <div className="battle-loading">
        <div className="loading-spinner"></div>
        <p>Loading battle...</p>
      </div>
    );
  }

  return (
    <div className="battle-container">
      <div className="battle-header">
        <div className="match-info">
          <h2>Battle: {problem.title}</h2>
          <span className="difficulty-badge">{problem.difficulty}</span>
        </div>
        
        <div className="timer">
          <span className="timer-icon">⏰</span>
          <span className="timer-text">{formatTime(timeLeft)}</span>
        </div>
        
        <div className="match-status">
          <span className={`status-indicator status-${matchStatus}`}></span>
          {matchStatus.toUpperCase()}
        </div>
      </div>

      <div className="battle-content">
        <div className="battle-side">
          <div className="player-header">
            <div className="player-info">
              <div className="player-avatar">{player1.username[0]}</div>
              <div>
                <div className="player-name">{player1.username}</div>
                <div className="player-stats">
                  <span>Score: {player1.score}</span>
                  <span>Submissions: {player1.submissions}</span>
                </div>
              </div>
            </div>
            <div className="player-status">
              <span className={`status-indicator status-${player1.status}`}></span>
              {player1.status}
            </div>
          </div>
          
          <div className="editor-container">
            <Editor
              height="100%"
              language="go"
              value={code}
              onChange={handleCodeChange}
              onMount={handleEditorDidMount}
              theme="vs-dark"
              options={{
                minimap: { enabled: false },
                fontSize: 14,
                lineNumbers: 'on',
                roundedSelection: false,
                scrollBeyondLastLine: false,
                automaticLayout: true,
              }}
            />
          </div>
          
          <div className="skill-cards">
            {skillCards.map(card => (
              <button
                key={card.id}
                className="skill-card"
                onClick={() => handleSkillCardUse(card)}
                title={card.description}
              >
                {card.name}
              </button>
            ))}
          </div>
          
          <div className="submit-section">
            <button
              className="submit-button"
              onClick={handleSubmit}
              disabled={isSubmitting || !code.trim()}
            >
              {isSubmitting ? 'Submitting...' : 'Submit Code'}
            </button>
          </div>
        </div>

        <div className="battle-side">
          <div className="player-header">
            <div className="player-info">
              <div className="player-avatar">{player2.username[0]}</div>
              <div>
                <div className="player-name">{player2.username}</div>
                <div className="player-stats">
                  <span>Score: {player2.score}</span>
                  <span>Submissions: {player2.submissions}</span>
                </div>
              </div>
            </div>
            <div className="player-status">
              <span className={`status-indicator status-${player2.status}`}></span>
              {player2.status}
            </div>
          </div>
          
          <div className="editor-container">
            <div className="opponent-editor">
              <div className="editor-placeholder">
                <p>Opponent's code is hidden</p>
                <p>Use skill cards to peek!</p>
              </div>
            </div>
          </div>
          
          <div className="skill-cards">
            <div className="opponent-cards">
              <span>Opponent's skill cards are hidden</span>
            </div>
          </div>
        </div>

        <div className="problem-panel">
          <div className="problem-header">
            <h3 className="problem-title">{problem.title}</h3>
            <span className="problem-difficulty">{problem.difficulty}</span>
          </div>
          
          <div className="problem-content">
            <div className="problem-description">
              {problem.description.split('\n').map((line, index) => (
                <p key={index}>{line}</p>
              ))}
            </div>
            
            <div className="test-cases">
              <h4>Test Cases</h4>
              {problem.testCases.map((testCase, index) => (
                <div key={index} className="test-case">
                  <div className="test-case-title">Test Case {index + 1}</div>
                  <div className="test-case-content">
                    <div><strong>Input:</strong> {testCase.input}</div>
                    <div><strong>Expected Output:</strong> {testCase.output}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Battle;

import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';
import '../custom.css';

export default function Dashboard() {
  const [technologies, setTechnologies] = useState([]);
  useEffect(() => {
    const fetchTechnologies = async () => {
      try {
        const response = await api.get('/api/technology');
        setTechnologies(response.data);
      } catch (error) {
        console.error('Error fetching technologies:', error);
      }
    };

    fetchTechnologies();
  }, []);

  return (
    <div className="dashboard-container">
      <div className="base-container">
        {technologies.map((tech) => (
          <div key={tech.id} className="base-card">
            <img src={tech.image_url} alt={tech.name} className="base-icon" />
            <Link to={`/technology/${tech.id}/components`} className="base-button">
              {tech.name}
            </Link>
            <p className="base-description">{tech.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
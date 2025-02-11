import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';

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
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {technologies.map((tech) => (
        <Link
          key={tech.id}
          to={`/technology/${tech.id}/components`}
          className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
        >
          <img
            src={tech.image_url}
            alt={tech.name}
            className="w-24 h-24 mx-auto mb-4 object-contain"
          />
          <h3 className="text-xl font-semibold text-center">{tech.name}</h3>
          <p className="text-gray-600 text-center mt-2">{tech.description}</p>
        </Link>
      ))}
    </div>
  );
}
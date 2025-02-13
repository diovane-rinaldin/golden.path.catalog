import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../services/api';

export default function ComponentList() {
  const { id } = useParams();
  const [components, setComponents] = useState([]);
  const [technology, setTechnology] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [techResponse, componentsResponse] = await Promise.all([
          api.get(`/api/technology/id/${id}`),
          api.get(`/api/component/technology/${id}`)
        ]);
        
        setTechnology(techResponse.data);
        setComponents(componentsResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [id]);

  if (!technology) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-container">
          <div className="base-container">
            {components.map((comp) => (
              <div key={comp.id} className="base-card">
                <img src={comp.image_url} alt={comp.name} className="base-icon" />
                <Link to={`/component/${comp.id}/services`} className="base-button">
                  {comp.name}
                </Link>
                <p className="base-description">{comp.description}</p>
              </div>
            ))}
          </div>
        </div>
  );
}
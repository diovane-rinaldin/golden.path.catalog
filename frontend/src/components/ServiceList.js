import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../services/api';

export default function ServiceList() {
  const { id } = useParams();
  const [services, setServices] = useState([]);
  const [component, setComponent] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [compResponse, servicesResponse] = await Promise.all([
          api.get(`/api/component/id/${id}`),
          api.get(`/api/service/component/${id}`)
        ]);
        
        setComponent(compResponse.data);
        setServices(servicesResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [id]);

  if (!component) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-container">
          <div className="base-container">
            {services.map((svc) => (
              <div key={svc.id} className="base-card">
                <img src={svc.image_url} alt={svc.name} className="base-icon" />
                <p>{svc.cloud_provider} - {svc.name}</p>
                <p className="base-description">{svc.description}</p>
              </div>
            ))}
          </div>
        </div>
  );
}
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
          api.get(`/api/component/${id}`),
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
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-bold">{component.name} Services</h2>
        <p className="text-gray-600">{component.description}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {services.map((service) => (
          <a
            key={service.id}
            href={service.service_cloud_url}
            target="_blank"
            rel="noopener noreferrer"
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
          >
            <img
              src={service.image_url}
              alt={service.name}
              className="w-24 h-24 mx-auto mb-4 object-contain"
            />
            <h3 className="text-xl font-semibold text-center">{service.name}</h3>
            <p className="text-gray-600 text-center mt-2">{service.description}</p>
            <div className="mt-4 text-center">
              <span className="inline-block bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm">
                {service.cloud_provider}
              </span>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
}
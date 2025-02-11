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
          api.get(`/api/technology/${id}`),
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
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-bold">{technology.name} Components</h2>
        <p className="text-gray-600">{technology.description}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {components.map((component) => (
          <Link
            key={component.id}
            to={`/component/${component.id}/services`}
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
          >
            <img
              src={component.image_url}
              alt={component.name}
              className="w-24 h-24 mx-auto mb-4 object-contain"
            />
            <h3 className="text-xl font-semibold text-center">{component.name}</h3>
            <p className="text-gray-600 text-center mt-2">{component.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
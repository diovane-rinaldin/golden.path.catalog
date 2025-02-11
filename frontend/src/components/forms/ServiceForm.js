import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';

export default function ServiceForm() {
  const navigate = useNavigate();
  const [components, setComponents] = useState([]);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    component_id: '',
    cloud_provider: '',
    service_cloud_url: '',
    image: null
  });

  useEffect(() => {
    const fetchComponents = async () => {
      try {
        const response = await api.get('/api/component');
        setComponents(response.data);
      } catch (error) {
        console.error('Error fetching components:', error);
      }
    };

    fetchComponents();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    try {
      const formDataObj = new FormData();
      if (formData.image) {
        formDataObj.append('image', formData.image);
        const uploadResponse = await api.post('/api/upload', formDataObj);
        
        const serviceData = {
          name: formData.name,
          description: formData.description,
          component_id: formData.component_id,
          cloud_provider: formData.cloud_provider,
          service_cloud_url: formData.service_cloud_url,
          image_url: uploadResponse.data.url
        };
        
        await api.post('/api/service', serviceData);
        navigate('/');
      }
    } catch (error) {
      console.error('Error creating service:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-lg mx-auto">
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Component
        </label>
        <select
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          value={formData.component_id}
          onChange={(e) => setFormData({ ...formData, component_id: e.target.value })}
          required
        >
          <option value="">Select a component</option>
          {components.map((comp) => (
            <option key={comp.id} value={comp.id}>
              {comp.name}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Service Name
        </label>
        <input
          type="text"
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          required
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Description
        </label>
        <textarea
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          required
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Cloud Provider
        </label>
        <select
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          value={formData.cloud_provider}
          onChange={(e) => setFormData({ ...formData, cloud_provider: e.target.value })}
          required
        >
          <option value="">Select a cloud provider</option>
          <option value="AWS">AWS</option>
          <option value="Azure">Azure</option>
          <option value="GCP">GCP</option>
        </select>
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Service Cloud URL
        </label>
        <input
          type="url"
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          value={formData.service_cloud_url}
          onChange={(e) => setFormData({ ...formData, service_cloud_url: e.target.value })}
          required
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Image
        </label>
        <input
          type="file"
          accept="image/*"
          onChange={(e) => setFormData({ ...formData, image: e.target.files[0] })}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          required
        />
      </div>

      <button
        type="submit"
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
      >
        Create Service
      </button>
    </form>
  );
}
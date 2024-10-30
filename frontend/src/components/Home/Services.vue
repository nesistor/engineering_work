<template>
    <div class="service-background">
      <div class="service-container">
        <!-- Nagłówek -->
        <header class="header">
          <h1>Microservices</h1>
          <p>You can test queries to several of my microservices here.</p>
        </header>
  
        <!-- Lista przycisków wyboru serwisu -->
        <div class="service-list">
          <ul>
            <li v-for="(service, index) in services" :key="index">
              <button @click="selectService(index)">
                {{ service.name }}
              </button>
            </li>
          </ul>
        </div>
  
        <!-- Widok dla wybranego serwisu -->
        <div class="request-response">
          <div class="request">
            <h2>Request</h2>
            <pre>{{ services[selectedService].request }}</pre>
            <button @click="sendRequest">Send</button>
          </div>
  
          <div class="response">
            <h2>Response</h2>
            <pre>{{ response }}</pre>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  
  <script>
   import '../../assets/css/Services.css';
  export default {
    name: 'TestServices',
    data() {
      return {
        selectedService: 0, // Początkowo wybrany serwis
        response: 'No response yet', // Pusty początkowy `response`
        services: [
          {
            name: 'Service 1',
            request: '{ "example": "Sample request for Service 1" }',
          },
          {
            name: 'Service 2',
            request: '{ "example": "Sample request for Service 2" }',
          },
          {
            name: 'Service 3',
            request: '{ "example": "Sample request for Service 3" }',
          },
          {
            name: 'Service 4',
            request: '{ "example": "Sample request for Service 4" }',
          },
          {
            name: 'Service 5',
            request: '{ "example": "Sample request for Service 5" }',
          },
        ],
      };
    },
    methods: {
      sendRequest() {
        // Symulujemy response, tutaj można podpiąć API
        this.response = `Response from ${this.services[this.selectedService].name}`;
      },
    },
  };
  </script>
  
<template>
    <div class="service-background">
      <div class="service-container">
        <!-- Nagłówek -->
        <header class="header">
          <h1>Wybór Serwisu</h1>
          <p>Wybierz serwis, aby zobaczyć szczegóły.</p>
        </header>
  
        <!-- Tło serwisu -->
        <div class="test-services">
          <!-- Formularz wyboru serwisu -->
          <div class="service-form">
            <label for="service-select">Select Service:</label>
            <select id="service-select" v-model="selectedService">
              <option v-for="(service, index) in services" :key="index" :value="index">
                {{ service.name }}
              </option>
            </select>
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
  
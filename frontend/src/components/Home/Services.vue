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

          <!-- Wyświetlany endpoint -->
          <p>Endpoint: {{ currentEndpoint.name }}</p>
          
          <!-- Nagłówek dla JSON -->
          <p>JSON:</p>
          <pre>{{ JSON.stringify(JSON.parse(currentEndpoint.request), null, 2) }}</pre>

          <!-- Wyświetlane nagłówki -->
          <p>Headers:</p>
          <pre>{{ JSON.stringify(currentEndpoint.headers, null, 2) }}</pre>

          <button @click="sendRequest">Send</button>
        </div>

        <div class="response">
          <h2>Response</h2>
          <pre>{{ formatResponse(response) }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import '../../assets/css/Services.css';
export default {
  name: 'TestServices',
  data() {
    return {
      selectedService: 0, // Początkowo wybrany serwis
      selectedEndpoint: 0, // Początkowo wybrany endpoint
      response: 'No response yet', // Pusty początkowy `response`
      services: [
        {
          name: 'auth-service',
          endpoints: [
            {
              name: '/api/auth/register',
              request: `{
                "username": "new_user",
                "email": "new_user@example.com",
                "password": "securePassword123",
                "firstName": "John",
                "lastName": "Doe",
                "dateOfBirth": "1990-01-01",
                "phone": "+1234567890"
              }`,
              headers: {
                "Content-Type": "application/json"
              }
            },
            {
              name: '/api/auth/login',
              request: `{
                "username": "new_user",
                "password": "securePassword123"
              }`,
              headers: {
                "Content-Type": "application/json"
              }
            },
          ],
        },
        {
          name: 'user-service',
          endpoints: [
            {
              name: '/api/user/get-profile',
              request: `{
                "userId": "12345"
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <token>"
              }
            },
            {
              name: '/api/user/update-profile',
              request: `{
                "userId": "12345",
                "firstName": "John",
                "lastName": "Doe",
                "email": "updated_email@example.com",
                "phone": "+0987654321"
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <token>"
              }
            },
          ],
        },
        {
          name: 'logger-service',
          endpoints: [
            {
              name: '/api/logger/log-event',
              request: `{
                "eventType": "USER_LOGIN",
                "userId": "12345",
                "timestamp": "2024-10-31T12:00:00Z"
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <token>"
              }
            },
            {
              name: '/api/logger/get-logs',
              request: `{
                "userId": "12345",
                "fromDate": "2024-01-01",
                "toDate": "2024-12-31"
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <token>"
              }
            },
          ],
        },
        {
          name: 'admin-service',
          endpoints: [
            {
              name: '/api/admin/add-user',
              request: `{
                "username": "admin_user",
                "email": "admin_user@example.com",
                "role": "administrator",
                "permissions": ["READ", "WRITE", "DELETE"]
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <admin_token>"
              }
            },
            {
              name: '/api/admin/remove-user',
              request: `{
                "userId": "67890",
                "reason": "Violation of terms"
              }`,
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer <admin_token>"
              }
            },
          ],
        },
      ],
    };
  },
  computed: {
    // Wybrany endpoint
    currentEndpoint() {
      return this.services[this.selectedService].endpoints[this.selectedEndpoint];
    },
  },
  methods: {
  selectService(index) {
    this.selectedService = index;
    this.selectedEndpoint = 0; // Reset selected endpoint on service change
  },
  async sendRequest() {
    // Parse the JSON request body from the current endpoint
    const requestBody = JSON.parse(this.currentEndpoint.request);

    try {
      // Send the request using Axios with real API interaction
      const res = await axios({
        method: 'post', // Use 'post' as the default or update per endpoint requirements
        url: `http://localhost:8080${this.currentEndpoint.name}`, // Replace with your API base URL
        headers: this.currentEndpoint.headers,
        data: requestBody
      });

      // Store the actual response data
      this.response = res.data;

    } catch (error) {
      // Handle errors gracefully and display error messages
      this.response = error.response ? error.response.data : 'Request failed. Please check the API URL and try again.';
    }

    // Automatically switch to the next endpoint, or reset if it's the last one
    if (this.selectedEndpoint < this.services[this.selectedService].endpoints.length - 1) {
      this.selectedEndpoint++;
    } else {
      this.selectedEndpoint = 0;
    }
  },
  formatResponse(response) {
    // Format the response in JSON with 2 spaces indentation for better readability
    return JSON.stringify(response, null, 2);
  },
},
};
</script>

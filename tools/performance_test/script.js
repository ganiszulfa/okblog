import http from 'k6/http';
import { sleep, check } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 7, 
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 10, 
      maxVUs: 50, 
    },
  },
};

export default function () {
  // Base URL from environment variable or default
  const baseUrl = __ENV.BASE_URL || 'http://localhost';
  const postCount = __ENV.POST_COUNT || 100;
  const tagCount = __ENV.TAG_COUNT || 10;
  const randomWordCount = __ENV.RANDOM_WORD_COUNT || 10;
  const pageCount = postCount / 10;
  
  // Generate random numbers for URLs
  const randomNumber = randomIntBetween(1, postCount);
  const tagNumber = randomIntBetween(1, tagCount);
  const randomPage = randomIntBetween(1, pageCount);
  const randomWord = `random_words_${randomIntBetween(1, randomWordCount)}`;
  
  // Array of URL patterns to choose from
  const urlPatterns = [
    baseUrl,
    `${baseUrl}/tag/tag${tagNumber}`,
    `${baseUrl}/test-post-${randomNumber}`,
    `${baseUrl}/?page=${randomPage}`,
    `${baseUrl}/search?t=${randomWord}`
  ];
  
  // Randomly select one URL pattern
  const randomIndex = randomIntBetween(0, urlPatterns.length - 1);
  const url = urlPatterns[randomIndex];
  
  // Send the request
  const response = http.get(url);
  
  // Check if the request was successful
  check(response, {
    'is status 200': (r) => r.status === 200,
    'transaction time < 500ms': (r) => r.timings.duration < 500,
  });
  
  console.log(`Visited: ${url}`);
  
}

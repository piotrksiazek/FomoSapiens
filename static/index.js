const priceChart = document.getElementById("priceChart").getContext("2d");
const input = document.getElementById("daysToFetch");
const chartButton = document.getElementById("chartButton");
const submissionButton = document.getElementById("submissionButton");
const date = document.getElementById("date");

const defaultNumverOfDaysToFetch = 180;

let usdPrice = [];
let sentiment = [];
let days = [];

const formatter = new Intl.DateTimeFormat("pl");

let chart = new Chart();

const newChart = () => {
  chart = new Chart(priceChart, {
    type: "line",
    data: {
      labels: days,
      datasets: [
        {
          label: "BTC price in USD",
          yAxisID: "price",
          data: usdPrice,
          fill: false,
          borderColor: "gold",
          tension: 0.1,
        },
        {
          label: "Reddit sentiment",
          yAxisID: "sentiment",
          data: sentiment,
          fill: false,
          borderColor: "red",
          tension: 0.1,
        },
      ],
    },
    options: {
      responsive: true, // Instruct chart js to respond nicely.
      maintainAspectRatio: false, // Add to prevent default behaviour of full-width/height
      plugins: {
        legend: {
          position: "bottom",
        },
      },

      scales: {
        sentiment: {
          type: "linear",
          display: true,
          position: "left",
          max: 50,
          min: 30,
        },
        price: {
          type: "linear",
          display: true,
          position: "right",

          // grid line settings
          grid: {
            drawOnChartArea: false, // only want the grid lines for one axis to show up
          },
        },
      },
    },
  });
};

const createChart = (howManyDays) => {
  axios
    .get(`getDays/${howManyDays}`)
    .then((response) => {
      response.data.forEach((day) => {
        usdPrice.push(day.realprice);
        sentiment.push(day.sentiment);
        days.push("");
      });
      newChart();
    })
    .catch((error) => console.log(error));
};

chartButton.addEventListener("click", () => {
  if (parseInt(input.value) < 0) return;
  usdPrice = [];
  sentiment = [];
  days = [];
  axios
    .get(`getDays/${input.value}`)
    .then((response) => {
      response.data.forEach((day) => {
        usdPrice.push(day.realprice);
        sentiment.push(day.sentiment);
        days.push("");
      });
      chart.destroy();
      newChart();
    })
    .catch((error) => console.log(error));
});

const redditCommentDiv = document.getElementById("redditComment");

submissionButton.addEventListener("click", () => {
  axios.get(`getTopSubmission/${date.value}`).then((response) => {
    console.log(response);
    redditCommentDiv.innerHTML = `
      <h3>${response.data.author}</h3>
      <p>${response.data.body}</p>
    `;
  });
});

createChart(defaultNumverOfDaysToFetch);

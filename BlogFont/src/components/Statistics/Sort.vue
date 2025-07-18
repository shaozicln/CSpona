<template>
    <div class="dashboard-container">
        <!-- 日历组件 -->
        <Calendar class="calendar-component" />
        
        <!-- 图表区域 -->
        <div class="chart-area">
            <button @click="changeChartType" id="chart-button">换一种统计方式？</button>
            <canvas id="chart" width="800" height="619"></canvas>
        </div>
    </div>
</template>

<script setup>
import Calendar from './Calendar.vue';

// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { onMounted } from 'vue';
import { Chart } from 'chart.js';
import { registerables } from 'chart.js';
Chart.register(...registerables);

let chart;
let chartType = 'bar';
Chart.defaults.scale.grid.color = 'rgba(0, 0, 0, 0.3)';
Chart.defaults.scale.ticks.color = 'rgba(0, 0, 0)';
Chart.defaults.scale.grid.display = false;
Chart.defaults.scale.ticks.display = true;

onMounted(() => {
    const articleCounts = JSON.parse(localStorage.getItem('articleCounts'));
    const labels = Object.keys(articleCounts).map(label => decodeURIComponent(label));
    const data = Object.values(articleCounts);

    const ctx = document.getElementById('chart').getContext('2d');
    chart = new Chart(ctx, {
        type: chartType,
        data: {
            labels,
            datasets: [{
                label: '文章数量',
                data,
                backgroundColor: [
                    'rgba(204,180, 44, 0.8)',
                    'rgba(220,239,233, 0.8)',
                    'rgba(164,214,211, 0.8)',
                    'rgba(116,164, 88, 0.8)',
                    'rgba( 62, 99, 66, 0.8)',
                ],
                gradient: {
                    backgroundColor: [
                        'rgba(204,180, 44, 0.8)',
                        'rgba(220,239,233, 0.8)',
                        'rgba(164,214,211, 0.8)',
                        'rgba(116,164, 88, 0.8)',
                        'rgba( 62, 99, 66, 0.8)',
                    ],
                    borderColor: 'rgba(0,0,0,0.8)',
                    borderWidth: 1
                },
                borderColor: chartType === 'line' ? 'black'
                    : [
                        // 'rgba(204,180, 44)',
                        // 'rgba(220,239,233)',
                        // 'rgba(164,214,211)',
                        // 'rgba(116,164, 88)',
                        // 'rgba( 62, 99, 66)',
                        'rgba(0,0,0,0.8)'
                    ],
                borderWidth: chartType === 'line' ? 3 : 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            title: {
                display: true,
                text: '分类数量统计'
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: '数量'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: '分类'
                    }
                }
            },
            legend: {
                display: true,
                labels: {
                    generateLabels: (chart) => {
                        const labels = chart.data.labels;
                        const colors = chart.data.datasets[0].backgroundColor;
                        return labels.map((label, index) => {
                            return {
                                text: label,
                                fillStyle: colors[index],
                                strokeStyle: colors[index],
                                lineWidth: 1,
                                hidden: false,
                                // 添加这行来使 legend 显示分类名称对应的颜色
                                // 对于饼图
                                pointStyle: chart.type === 'pie' ? 'circle' : 'rect'
                            };
                        });
                    }
                }
            },
            animation: {
                duration: 1500,
                easing: 'easeInOutQuad'
            }
        }
    });
});

function changeChartType() {
    chart.destroy();
    chartType = chartType === 'bar' ? 'line' : chartType === 'line' ? 'pie' : chartType === 'pie' ? 'bubble' : chartType === 'bubble' ? 'doughnut' : chartType === 'doughnut' ? 'polarArea' : chartType === 'polarArea' ? 'radar' : chartType === 'radar' ? 'scatter' : 'bar';
    const articleCounts = JSON.parse(localStorage.getItem('articleCounts'));
    const labels = Object.keys(articleCounts).map(label => decodeURIComponent(label));
    const data = Object.values(articleCounts);
    const ctx = document.getElementById('chart').getContext('2d');
    chart = new Chart(ctx, {
        type: chartType,
        data: {
            labels,
            datasets: [{
                label: '文章数量',
                data,
                backgroundColor: [
                    'rgba(204,180, 44, 0.8)',
                    'rgba(220,239,233, 0.8)',
                    'rgba(164,214,211, 0.8)',
                    'rgba(116,164, 88, 0.8)',
                    'rgba( 62, 99, 66, 0.8)',
                ],
                borderColor: chartType === 'line' ? 'black'
                    : [
                        // 'rgba(204,180, 44)',
                        // 'rgba(220,239,233)',
                        // 'rgba(164,214,211)',
                        // 'rgba(116,164, 88)',
                        // 'rgba( 62, 99, 66)',
                        'rgba(0,0,0,0.8)'
                    ],
                borderWidth: chartType === 'line' ? 3 : 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            title: {
                display: true,
                text: '分类数量统计'
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: '数量'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: '分类'
                    }
                }
            },
            legend: {
                display: true,
                labels: {
                    generateLabels: (chart) => {
                        const labels = chart.data.labels;
                        const colors = chart.data.datasets[0].backgroundColor;
                        return labels.map((label, index) => {
                            return {
                                text: label,
                                fillStyle: colors[index],
                                strokeStyle: colors[index],
                                lineWidth: 1,
                                hidden: false
                            };
                        });
                    }
                }
            },
            animation: {
                duration: 1500,
                easing: 'easeInOutQuad'
            }
        }
    });
}

import { onBeforeRouteLeave } from 'vue-router';
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});
</script>

<style scoped>
.dashboard-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 20px;
}

.calendar-component {
    background-color: rgba(255, 255, 255, 0.3);
    border-radius: 8px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    padding: 15px;
}

.chart-area {
    position: relative;
    background-color: rgba(255, 255, 255, 0.3);
    border-radius: 8px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    padding: 60px;
}

#chart-button {
    position: absolute;
    top: 20px;
    left: 20px;
    border-radius: 20px;
    width: 170px;
    height: 40px;
    background-color: rgba(97, 138, 184, 0.684);
    font-size: 17px;
    z-index: 1;
    border: none;
    cursor: pointer;
    transition: background-color 0.3s;
}

#chart-button:hover {
    background-color: rgba(77, 118, 164, 0.784);
}
</style>
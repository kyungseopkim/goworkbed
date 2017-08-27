$(function() {
    $.getJSON("http://localhost:1323/operation/bydoctor?dbname=testdb", function(data) {
        opsdata = []
        doctors = []
        $.each(data, function(index, val) {
            doctors.push(val['의사'])
            opsdata.push(val['수술'])
        });
        // container
        Highcharts.chart('container', {
            chart: {
                type: 'column'
            },
            title: {
                text: 'Operation By Doctors'
            },
            xAxis: {
                categories: doctors,
                scrollbar: {
                    enabled: true
                },
                ticklength: 0,
                title: {
                    text: null
                },
                min: 0,
                max: 10
            },
            yAxis: {
                title: {
                    text: 'Operation'
                }
            },
            series: [{
                name: '수술',
                data: opsdata
            }],
            plotOptions: {
                bar: {
                    dataLabels: {
                        enabled: true
                    }
                }
            },

            legend: {
                enabled: false
            },
            credits: {
                enabled: false
            }
        });
        // container
        Highcharts.chart('container2', {
            chart: {
                type: 'column'
            },
            title: {
                text: 'Operation By Doctors'
            },
            xAxis: {
                categories: doctors,
                ticklength: 0,
                title: {
                    text: null
                }
            },
            yAxis: {
                title: {
                    text: 'Operation'
                }
            },
            series: [{
                name: '수술',
                data: opsdata
            }],
            plotOptions: {
                bar: {
                    dataLabels: {
                        enabled: true
                    }
                }
            },

            legend: {
                enabled: false
            },
            credits: {
                enabled: false
            }
        });

    });
});

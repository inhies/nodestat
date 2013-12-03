//
// api.js
//

$(document).ready(function() {
    
    // Originally we want to get the time
    getInfo();
    
    // Now every second we want to increment the last updated
    // thing. Every 5 seconds we want to update the info.
    setInterval(lastUpdated, 1000);

    // We want to bind some click events for core, angel, and peers
    // now so we can switch between them.
    $('#angel').bind('click', function() {
	$('.core, .peers').css('display', 'none');
	$('.angel').css('display', 'block');
	$('#angel').addClass('active');
	$('#core, #peers').removeClass('active');
    });
    $('#core').bind('click', function() {
	$('.angel, .peers').css('display', 'none');
	$('.core').css('display', 'block');
	$('#core').addClass('active');
	$('#angel, #peers').removeClass('active');
    });
    $('#peers').bind('click', function() {
	$('.core, .angel').css('display', 'none');
	$('.peers').css('display', 'block');
	$('#peers').addClass('active');
	$('#angel, #core').removeClass('active');
    });
    
});

function lastUpdated() {
    var time = $('#update').html().split(' ');
    if (time[0] == 5) {
	getInfo();
	time[0] = 0;
    } else {
	time[0]++;
    }
    $('#update').html(time.join(' '));
}

function getInfo() {
    $.ajax ({
	url: '/all',
	type: 'GET',
	dataType: 'json',
	statusCode: {
            200: function(res) {
                $('.header').removeClass('failure').addClass('success');
                $('.header').html('All Systems Operational');
		
                $('.core .uptime').html('<div class="text">Uptime</div>'+res.Node.Core.Uptime);
                $('.core .cpu').html('<div class="text">CPU Usage</div>'+res.Node.Core.PercentCPU+'%');
                $('.core .mem').html('<div class="text">Memory Usage</div>'+res.Node.Core.PercentMemory+'%');
                $('.angel .uptime').html('<div class="text">Uptime</div>'+res.Node.Angel.Uptime);
                $('.angel .cpu').html('<div class="text">CPU Usage</div>'+res.Node.Angel.PercentCPU+'%');
                $('.angel .mem').html('<div class="text">Memory Usage</div>'+res.Node.Angel.PercentMemory+'%');
		
                // Peers                                                                                                                                                         
                var peers = 0;
                $.each(res.Peers, function(key, val) {
                    peers++;
                });
                $('.peers').html('<pre>'+JSON.stringify(res.Peers, null, 4)+'</pre>');
                $('#peers').html('Peers ('+peers+')');
            },
            404: function() {
                $('.header').removeClass('success').addClass('failure');
                $('.header').html('Unable to connect to server API!');
            },
            401: function() {
                $('.header').removeClass('success').addClass('failure');
                $('.header').html('You are not authorized to view this information!');
            },
            500: function() {
                $('.header').removeClass('success').addClass('failure');
                $('.header').html('Server error!');
            },
            502: function() {
                $('.header').removeClass('success').addClass('failure');
                $('.header').html('Unable to connect to server API!');
            },
            503: function() {
                $('.header').removeClass('success').addClass('failure');
                $('.header').html('Cjdns is not running!');
            }
        }
    });
    
}

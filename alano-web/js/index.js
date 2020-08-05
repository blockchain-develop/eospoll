(function ($) {

    // 投注金额
    $('#Half').click(function () {
        var money = 1
        $("#money").val(money);
        $('#may_get_money').val(money);

    })
    $('#Double').click(function () {
        var money = 2
        $("#money").val(money);
        $('#may_get_money').val(money);

    })
    $('#Max').click(function () {
        var money = 4
        $("#money").val(money);
        $('#may_get_money').val(money);

    })
    
    $("#money").bind("input propertychange", function() {
    	var selected_number = $("#selected-number").val();
    	selected_number = parseInt(selected_number)
    	var money = $("#money").val();
    	money = parseInt(money)
    	calReadyBet(selected_number, money)
    })

    $("#money").change(function () {
    	{
	        var money = $("#money").val();
	        $('#may_get_money').val(money);
    	}
        
        //
    	{
	    	var selected_number = $("#selected-number").val();
	    	selected_number = parseInt(selected_number)
	    	var money = $("#money").val();
	    	money = parseInt(money)
	    	calReadyBet(selected_number, money)
    	}
    })
    
    $("#may_get_money").change(function () {
        var money = $("#may_get_money").val();
        $('#money').val(money);
    })
    
    $("#selected-number").change(function() {
    	var selected_number = $("#selected-number").val();
    	selected_number = parseInt(selected_number)
    	var money = $("#money").val();
    	money = parseInt(money)
    	calReadyBet(selected_number, money)
    })

    // 进度条数字获取  $bt.html( parseInt(left / 6.5));
    var $box = $('#box');
    var $result_bt = $('#result_bt'); //结果数字
    var $bg = $('#bg');
    var $bgcolor = $('#bgcolor');
    var $bgcolor_neg = $('#bgcolor-neg');
    var $btn = $('#dice-slider-button-i'); //选择的数
    var $slider_number = $('#dice-slider-button-em'); //选择的数

    var $text = $('#text');
    var statu = false;
    var ox = 0;
    var lx = 0;
    var left = 0;
    var bgleft = 0;
    var proxyer1 = "";
    var proxyer2 = "";

    var remain_draw_times = 0;

    $(".lottery_button").attr("disabled", true);

    $btn.on("mousedown", function (e) {
        lx = parseInt($btn.css("marginLeft")) + $bg.width() / 2;
        ox = e.pageX - lx;
        statu = true;
    });
    $btn.on("touchstart", function (e) {
        lx = parseInt($btn.css("marginLeft")) + $bg.width() / 2;
        ox = e.touches[0].clientX - lx;
        statu = true;
    });
    $(document).on("mouseup touchend", function (e) {
        statu = false;
    });


    $(document).on("mousemove", function (e) {
        max = $("#bg").width() - 30;
        min = max / 100;

        if (statu) {
            onChangeBet(e.pageX);
        }
    });
    $(document).on("touchmove", function (e) {
        max = $("#bg").width() - 30;
        min = max / 100;

        if (statu) {
            onChangeBet(e.touches[0].clientX);
        }
    });
    $bg.click(function (e) {
        max = $("#bg").width() - 30;
        min = max / 100;
        ox = $btn.offset().left - parseInt($btn.css("marginLeft")) - $bg.width() / 2;
        if (!statu) {
            onChangeBet(e.pageX);
        }
    });

    function onChangeBet(x) {
        left = x - ox;

        //max = max-min*3
        if (left < min * 2) {
            left = min * 2;
        }
        if (left > max - min * 4) {
            left = max - min * 4;
        }

        width = left - $bg.width() / 2;
        $btn.css('marginLeft', width);
        $slider_number.css('marginLeft', width - 15)

        $bgcolor.width(left + 15);
        $bgcolor_neg.width(max - left + 15);
        ratio = left * 100 / max;
        //$btn.html(parseInt(ratio));
        $slider_number.text(parseInt(ratio));

        $('#myNumber').html(parseInt(ratio));
        $('#myNumberValue').val(parseInt(ratio));
        var odds = Number(98.5 / (parseInt(ratio) - 1)).toFixed(2)
        $('#odds').html(odds + 'x'); //赔率计算
        $('#percent').html(parseInt((parseInt(ratio - 1) / 100) * 100) + '%'); //中奖概率计算
        $('#may_get_money').val(Number((odds * $("#money").val()).toFixed(3))); //可能获得的奖金
    }

    // 定义玩法函数
    var account = null;
    var eoss = null;
    var requiredFields = null;
    var tpAccount = null;
    var balanceEos = 0;
    var balanceBetDice = 0;
    var betContract = "alanogoddess"
    var bugContract = "alanogoddess"

    var inviteCode = "";
    var playType = 'eos'
    	
    var btime = 0;
    	
    var connected = false;
    
    var numberState = null;
    
    var betState = null;

    var urls = window.location.href.split('ref=')
    if (urls.length > 1 && urls[1].length >=3 && urls[1].length <=13) {
        inviteCode = urls[1]
    }
    console.log(inviteCode)
    $("#inviteLink").val(document.location.origin + "/?ref=");

    var a = document.getElementById("#result");

    var hideLoading = function () {
        $("#loading").modal("hide");
    }
    
    var init_scatter = function () {
        if (account != null) return;
        if (eoss != null) return;
        if (tpAccount != null) return;

        var that = this
        var status = false
        // 判断是否登录
        if (1) {
            var checkCount = 0
            var checkInterval = setInterval(function () {
                console.log(checkCount)
                if (checkCount > 20) {
                    clearInterval(checkInterval)
                    hideLoading()

                    showAlert('加载Scatter失败若未安装请先安装', false)
                    return
                }
                if (typeof (scatter) == "undefined") {
                    checkCount++
                    return
                } else {
                    if (status == true) {
                        return
                    }
                    status = true
                    clearInterval(checkInterval)

                    eoss = scatter.eos(network, Eos, {});
                    setInterval(function () {
                        getBetCurrentId()
                    }, 1000)
                    
                    
                    setInterval(function () {
                        getHistoryBetState()
                    }, 1000)
                    
                    $("#alertmsg").addClass("result_animation");
                    setTimeout('$("#alertmsg").removeClass("result_animation");', 4000);

                    scatter.getIdentity({
                            accounts: [{
                                chainId: network.chainId,
                                blockchain: network.blockchain
                            }]
                        })
                        .then(identity => {
                            setIdentity(identity);

                            showSuccess('登录成功')
                            if (account) {
                                $("#login").hide();
                                $('.nickname').html(account.name);
                                $("#play").text("我要下注")

                                $("#inviteLink").val(document.location.origin + "/?ref=" + account.name);
                            }

                            //get_bonuspool();

                            setTimeout(function () {
                                get_current_balance();
                            }, 1000)

                            setTimeout(function () {
                                get_cpu();
                                //get_invite_info();
                            }, 5000)
                            
                            setTimeout(function () {
                            	getMyInvite()
                            }, 1000)
                        })
                        .catch(err => {
                            console.log(err)
                            //alert( "Scatter 初始化失败.", err );

                            showAlert('Scatter 初始化失败, 请刷新重试', true)
                        });

                    hideLoading()
                }

                checkCount++
            }, 500)

        } else {
            //移动端
            tpConnected = tp.isConnected();
            if (tpConnected) {
                //test
                // app.tpBalance();
                tp.getWalletList("eos").then(function (data) {
                    tpAccount = data.wallets.eos[0];
                });
            } else {
                alert("请下载TokenPocket") //待完善
                hideLoading()
            }
        }
    };
    var setIdentity = function (identity) {

        account = identity.accounts.find(acc => acc.blockchain === 'eos');

        console.log("account", account)
        eoss = scatter.eos(network, Eos, {});
        requiredFields = {
            accounts: [network]
        };
        //get_current_balance();
    };

    //获取账户eos余额
    var get_current_balance = function () {
        eoss.getCurrencyBalance('eosio.token', account.name).then(function (resp) {
            console.log("get_current_balance", resp);
            balanceEos = resp[0]
            $('#balanceEos').text(balanceEos);
        });
    };
    
    var get_local_balance = function() {
        eoss.getTableRows({
            code: betContract, //EOS_CONFIG.contractName,
            scope: betContract, //.contractName,
            table: "betaccount",
            index_position:"2",
            key_type:"name",
            lower_bound:account.name
            limit: 1,
            json: true
        }).then(data => {
            console.log("get_local_balance", data)
            $("#bet-rank").html(html)
        }).catch(e => {
            console.error("offerbet ", e);
        });	
    }

    var get_cpu = function () {
        eoss.getAccount({
            account_name: account.name
        }).then(data => {
            cp = data.cpu_limit.max == 0 ? 1 : data.cpu_limit.used / data.cpu_limit.max
            np = data.net_limit.max == 0 ? 1 : data.net_limit.used / data.net_limit.max

            cp = cp > 1 ? 1 : cp;
            np = np > 1 ? 1 : np;
            net.animate(np); // Number from 0.0 to 1.0
            cpu.animate(cp);
        }).catch(e => {
            console.error("getAccout ", e);
        });
    }


    var roll_by_scatter = function () {
        var money = $("#money").val()
        money = parseInt(money * 10000) / 10000
        money = money.toFixed(4)
        
        var number = $("#selected-number").val()
        number = parseInt(number)
        {
            if (money < 1 || money > parseFloat(balanceEos)) {

                hideLoading()
                showAlert('EOS余额不足')
                return
            }
            money += " EOS"
        }
        
        {
        	if (number <= 0 || number > 12) {
        		showAlert("没有选择号码")
        		return;
        	}
        	
        	if(number == 12 && numberState[12].bets > 0) {
        		showAlert("限注号码已满")
        		return;	
        	}
        	
        	if(number == 11 && numberState[11].bets > 0) {
        		showAlert("限注号码已满")
        		return;	
        	}
        }

        eoss.contract(code, {
                accouts: [network]
            }).then(contract => {

                //var meno = $('#myNumberValue').val()
                //meno += ' ' + bet_id+ ' ' + inviteCode
            	var meno = ("0" + number).substr(-2)
            	meno = meno + inviteCode
                console.log(account.name)
                console.log(contract_name)
                console.log(meno)
                console.log(money)
                contract.transfer(account.name, contract_name, money, meno, {
                        authorization: [account.name + '@active']
                    })
                    //eoss.transfer(account.name, contract, money, $( '#myNumber' ).html()) // 抵押 换成‘stake’  赎回unstaketoken
                    .then((resp) => {

                        /*
                        var inline_traces = resp.processed.action_traces[0].inline_traces
                        var i = inline_traces.length - 1
                        var roll = inline_traces[i].act.data.result.random_roll
                        var payout = inline_traces[i].act.data.result.payout
                        */
                    	hideLoading()
                        get_current_balance();
                    })
                    .catch((err) => {
                        hideLoading()
                        //console.log( "err ", err, JSON.stringify( err ) );
                        if (err && err.message) {
                            showAlert('执行失败' + err.message, false)
                        } else {
                            var errJson = JSON.parse(err)

                            showAlert('执行失败' + errJson.error.name, false)
                        }

                    });
            })
            .catch((err) => {
                hideLoading()
                console.log("err ", err, JSON.stringify(err));

            });
    };

    var deposity = function (amount) {
        money = parseFloat(amount)
        money = money.toFixed(4)
        money += " EOS"
        var contract_name = betContract;
        var code = "eosio.token"
        var memo = "bet deposity"

        eoss.contract(code, {
            	accouts: [network]
            }).then(contract => {
                console.log(account.name)
                console.log(contract_name)
                console.log(money)
                console.log(memo)
                contract.transfer(account.name, contract_name, money, memo, {
                        authorization: [account.name + '@active']
                    }).then((resp) => {
                    	hideLoading()
                        get_current_balance();
                    }).catch((err) => {
                        hideLoading()
                        if (err && err.message) {
                            showAlert('执行失败' + err.message, false)
                        } else {
                            var errJson = JSON.parse(err)
                            showAlert('执行失败' + errJson.error.name, false)
                        }
                    });
            }).catch((err) => {
                hideLoading()
                console.log("err ", err, JSON.stringify(err));
            });
    };

    var offerbet = function (amount, number) {
        money = parseFloat(amount)
        money = money.toFixed(4)
        money += " EOS"
        var code = betContract
        eoss.contract(code, {
                accouts: [network]
            }).then(contract => {
                console.log(account.name)
                console.log(money)
                console.log(number)
                contract.offerbet(account.name, money, number, {
                        authorization: [account.name + '@active']
                    }).then((resp) => {
                    	hideLoading()
                    	get_current_balance();
                    }).catch((err) => {
                        hideLoading()
                        if (err && err.message) {
                            showAlert('执行失败' + err.message, false)
                        } else {
                            var errJson = JSON.parse(err)
                            showAlert('执行失败' + errJson.error.name, false)
                        }
                    });
            }).catch((err) => {
                hideLoading()
                console.log("err ", err, JSON.stringify(err));
            });
    };

    var withdraw = function (amount) {
        money = parseFloat(amount)
        money = money.toFixed(4)
        money += " EOS"
        var code = betContract
        eoss.contract(code, {
                accouts: [network]
            }).then(contract => {
                contract.withdraw({
                        authorization: [account.name + '@active']
                    }).then((resp) => {
                    	hideLoading()
                        get_current_balance();
                    }).catch((err) => {
                        hideLoading()
                        if (err && err.message) {
                            showAlert('执行失败' + err.message, false)
                        } else {
                            var errJson = JSON.parse(err)
                            showAlert('执行失败' + errJson.error.name, false)
                        }
                    });
            }).catch((err) => {
                hideLoading()
                console.log("err ", err, JSON.stringify(err));
            });
    };


    var bets = [];
    var currentId = 0;

    var getBetCurrentId = function () {
        eoss.getTableRows({
            code: betContract, //EOS_CONFIG.contractName,
            scope: betContract, //.contractName,
            table: "globalindex",
            json: true
        }).then(data => {
            console.log("globalindex ", data)

            currentId = data.rows[0].gindex

            //currentId = currentId > 20 ? currentId - 20 : -1;
            //getBetList()
            //getMyBetList()

            setTimeout(function () {
                getCurrencyBetList()
            }, 1000)
            
            
            setTimeout(function() {
            	getCurrencyBetState()
            }, 1000)
            
            setTimeout(function() {
            	getCurrencyBetLedger()
            }, 1000)
            
            setTimeout(function() {
            	getCurrencyNumberState()
            }, 1000)
        }).catch(e => {
            console.error("getTableRows ", e);
        });
    }

    var lottery_draw_time_interval;
    var last_draw_time;
    var getCurrencyBetList = function () {
        eoss.getTableRows({
            code: betContract, //EOS_CONFIG.contractName,
            scope: betContract, //.contractName,
            table: "betoffer",
            limit: 1000,
            json: true
        }).then(data => {
            console.log("betoffer", data)

            var html = '';
            var rows = data.rows;

            for (var i = rows.length - 1;i >= 0;i --)
            {
                var row = rows[i]
                html += '<li class="items top' + i + '">' +
                    '<span>' + row.user + '</span>' +
                    '<span>' + row.number + '</span>' +
                    '<span>' +
                    '<em style="font-weight: 600;">' + (row.quantity) + '</em>' +
                    '<em style="font-size: 12px; padding-left: 4px;"></em>' +
                    '</span>' +
                    '</li>';
            }
            $("#bet-rank").html(html)
        }).catch(e => {
            console.error("offerbet ", e);
        });
    }
    
    var getCurrencyNumberState = function () {
    	console.log("getCurrencyNumberState")
    	try {
    		var url = "http://103.45.157.202:7780/numberpool/all";
	    	$.getJSON(url, function(data) {
	    		console.log("getCurrencyNumberState callback")
	    		numberState = data
	    		
	    		var html = ""
	    		for(var i = 1;i < data.length;i ++) {
	    			var row = data[i]	    			
	                html += '<li class="items top' + i + '">' +
                    '<span>' + row.number + '</span>' +
                    '<span>' + row.bets + '</span>' +
                    '<span>' + row.precent + '</span>' +
                    '<span>' +
                    '<em style="font-weight: 600;">' + (row.prize) + '</em>' +
                    '<em style="font-size: 12px; padding-left: 4px;"></em>' +
                    '</span>' +
                    '</li>';
	    		}
	    		$("#bet-number-state").html(html)
	    	});
    	}
    	catch (ex) {
    		console.log(ex.message)
    	}
    }
    
    var getCurrencyBetLedger = function () {
        eoss.getTableRows({
            code: betContract, //EOS_CONFIG.contractName,
            scope: betContract, //.contractName,
            table: "innerledger",
            limit: 1000,
            json: true
        }).then(data => {
            console.log("innerledger", data)

            /*
            var bonus_total = 0;
            var bonus_me = 0;
            
            var rows = data.rows;
            for (var i = 0;i < rows.length;i ++)
            {
                var row = rows[i]
                if(row.user == account.name) {
                	bonus_me = row.quantity;
                	break;
                }
            }
            $("#bonus-me").text(bonus_me)
            */
        }).catch(e => {
            console.error("innerledger ", e);
        });
    }
    
    var getMyInvite = function () {
    	console.log("getMyInvite")
    	try {
    		var url = "http://103.45.157.202:7780/proxyermembers/";
    		url += account.name;
	    	$.getJSON(url, function(data){
	    		console.log("getMyInvite callback")
	    		var html = ""
	    		for(var i = 0;i < data.length;i ++) {
	    			var row = data[i]
	    			html += '<li id="li' + (i + 1) + '" class="normal-li" onclick="inviteselected(this,' + (i+1) + ',' + row + ')">' +
	    				'<span>' + row + '</span>' +
	    				'</li>';
	    		}
	    		$("#myinvite1-list").html(html)
	    		
	    		//
	    		if(data.length > 0) {
	    			proxyer1 = data[0];
	    		} else {
	    			proxyer1 = "";
	    		}
	    		
	    		//
	    		getMyInvite1();
	    	});
    	}
    	catch (ex) {
    		console.log(ex.message)
    	}
    }

    var getMyInvite1 = function () {
    	console.log("getMyInvite1")
    	if(proxyer1 == "") {
    		return;
    	}
    	
    	try {
    		var url = "http://103.45.157.202:7780/proxyermembers/";
    		url += proxyer1;
	    	$.getJSON(url, function(data){
	    		console.log("getMyInvite1 callback")
	    		var html = ""
	    		for(var i = 0;i < data.length;i ++) {
	    			var row = data[i]
	    			html += '<li id="li' + (i + 1) + '" class="normal-li" onclick="inviteselected(this,' + (i+1) + ',' + row + ')">' +
	    				'<span>' + row + '</span>' +
	    				'</li>';
	    		}
	    		$("#myinvite2-list").html(html)
	    		
	    		//
	    		if(data.length > 0) {
	    			proyer2 = data[0];
	    		} else {
	    			proxyer2 = "";
	    		}
	    	});
    	}
    	catch (ex) {
    		console.log(ex.message)
    	}
    }
    
    function inviteselected(obj, i, name) {
    	if(name == "") {
    		return
    	}
    	
    	proxyer1 = name;
    	getMyInvite1()
    }
    
    var getCurrencyBetState = function() {
        eoss.getTableRows({
            code: betContract, //EOS_CONFIG.contractName,
            scope: betContract, //.contractName,
            table: "betstate",
            limit: 1,
            lower_bound:currentId,
            json: true
        }).then(data => {
            console.log("betstate ", data)
            if(data.rows.length != 1) {
            	return;
            }
            
            betState = data.rows[0];
  
            var btime = new Date(betState.btime/1000);
            btime.setMinutes(parseInt(btime.getMinutes()/10)*10);
            
            $("#myNumber").text(btime.Format("yyyyMMddHHmm").slice(2));
            //$("#myNumber").text(row.round);
            $("#odds").text(betState.users);
            $("#percent").text(betState.quantity);
            $("#all_pool").text(betState.apoolamount);
            $("#bonus-total").text(betState.prizeamount);
            
        }).catch(e => {
            console.error("betstate ", e);
        });
    	
    }
    
    var getHistoryBetState = function() {
    	console.log("getHistoryBetState")
    	try {
	    	$.getJSON("http://103.45.157.202:7780/historybetstate/all", function(data){
	    		console.log("getHistoryBetState callback")
	    		var html = ""
	    		for(var i = 0;i < data.length;i ++) {
	    			var row = data[i]
	    			
	                var btime = new Date(row.btime_a/1000);
	                btime.setMinutes(parseInt(btime.getMinutes()/10)*10);
	                
	                html += '<tr class="xxx">' +
	                '<td>' + btime.Format("yyyyMMddHHmm").slice(2) + '</td>' +
	                '<td>' + row.otime + '</td>' +
	                '<td>' + row.users + '</td>' +
	                '<td>' + row.bets + '</td>' +
	                '<td>' + row.quantity + '</td>' +
	                '<td>' + row.inviteamount + '</td>' +
	                '<td>' + row.result + '</td>' +
	                '<td>' + 0 + '</td>' +
	                '<td>' + row.prizeamount + '</td>' +
	                '</tr>';
	    		}
	    		$("#history-betstate-list").html(html)
	    	});
    	}
    	catch (ex) {
    		console.log(ex.Message())
    	}
    }
   
    var calReadyBet = function(number,bets) {
    	console.log("calReadyBet",number)
    	console.log("calReadyBet",bets)
    	if(number <= 0 || number > 12 || numberState == null) {
    		$("#ready-number").text("N");
    		$("#ready-bets").text("N");
    		$("#ready-precent").text("N");
    		$("#ready-prize").text("N");
    		$("#ready-all-prize").text("N");
    		$("#number-state").text("无");
    	} else {
    		//
    		var precent = (numberState[number].bets + bets)/(numberState[0].bets + bets) * 100
    		precent = precent.toFixed(2)
    		console.log("precent",precent)
    		
    		var t_pool = parseFloat(betState.apoolamount) * 10000
    		var prize = (numberState[0].prize_a + bets*10000 + t_pool/2)/(numberState[number].bets + bets)
    		prize = prize/10000
    		prize = prize.toFixed(4)
    		console.log("prize",prize)
    		
    		var prize_all = prize * bets
    		prize_all = prize_all.toFixed(4)
    		console.log("prize_all",prize_all)
    		
    		if(number == 11 || number == 12) {
    			$("#number-state").text("限1位用户");
    		}
    		if(number == 9 || number == 10) {
    			$("#number-state").text("限2位用户");
    		}
    		if(number == 7 || number == 8) {
    			$("#number-state").text("限4位用户");
    		}
    		if(number == 5 || number == 6) {
    			$("#number-state").text("限8位用户");
    		}
    		if(number == 3 || number == 4) {
    			$("#number-state").text("限16位用户");
    		}
    		if(number == 1 || number == 2) {
    			$("#number-state").text("限32位用户");
    		}
    		
    		$("#ready-number").text(number);
    		$("#ready-bets").text(bets);
    		$("#ready-precent").text(precent + "%");
    		$("#ready-prize").text(prize + " EOS");
    		$("#ready-all-prize").text(prize_all + " EOS");
    	}
    }


    //保留4位小数并格式化输出（不足的部分补0）
    var fomatFloat = function (value, n) {
        var f = Math.round(value * Math.pow(10, n)) / Math.pow(10, n);
        var s = f.toString();
        var rs = s.indexOf('.');
        if (rs < 0) {
            s += '.';
        }
        for (var i = s.length - s.indexOf('.'); i <= n; i++) {
            s += "0";
        }
        return s;
    }

    $(".lottery_button").click(function () {

        $(this).attr("disabled", true)

        console.log("lottery_button");
        var that = this

        eoss.contract(betContract, {
            accouts: [network]
        }).then(contract => {

            contract.draw(account.name, {
                    authorization: [account.name + '@active']
                })
                .then((resp) => {
                    console.log(resp)

                    var inline_traces = resp.processed.action_traces[0].inline_traces
                    var i = inline_traces.length - 1
                    var luckey_num = inline_traces[i].act.data.luckey_num
                    var payout = inline_traces[i].act.data.payout
                    //console.log(payout)

                    //payout = parseFloat(payout)
                    //payout = payout/10000
                    showSuccess('抽奖成功! 幸运数字' + luckey_num + " 奖励" + payout);
                    $(that).attr("disabled", false)

                    getBetRanks()

                }).catch(err => {

                    if (err && err.message) {
                        showAlert('执行失败' + err.message, false)
                    } else {
                        var errJson = JSON.parse(err)

                        showAlert('抽奖失敗' + errJson.error.name);
                    }

                    $(that).attr("disabled", false)
                });
        })
    })

    var bonus_eos_amount, bonus_token_amount, bonus_stake_amount;



    var release_time_interval;
    var release_time;
    var showSuccess = function (msg) {
        $(".modal.alert-msg").modal("show");

        var html = '<div  class="alert alert-success" style="padding:20px 0;">' +
            msg +
            '</div>';
        $("#alert-msg").html(html)
        setTimeout(function () {
            $("#alert-modal").modal("hide");
        }, 1500)
    }

    var showAlert = function (msg, isShow) {
        var html = '<div  class="alert alert-warning" style="padding:20px 0;">';
        if (isShow == undefined || !isShow) {
            html += '<a href="#" class="close" data-dismiss="alert">x</a>';
        }

        html += msg + '</div>';

        $("#alert-msg").html(html)
        $("#alert-modal").modal("show");
        if (isShow == undefined || !isShow) {
            setTimeout(function () {
                $("#alert-modal").modal("hide");
                hideLoading()
            }, 2500)
        }
    }

    var checkLogin = function () {
        if (account == null) {
            $("#play").text("登录中...")
            //$(".modal.loading").modal("show");
            init_scatter();
            return true
        }
        return false;
    }
    // play
    $('#play').click(function () {
        $("#loading").modal("show");
        //init_scatter();

        if (checkLogin()) {
            return
        }
        roll_by_scatter();
    })

    // play
    $('#login').click(function () {
        if (account != null) {
        	ScatterJS.scatter.forgetIdentity();
        	return;
        } else {
	        $("#loading").modal("show");
	        init_scatter();
	    }
    })
    
    
    var connect_scatter = function () {
        setTimeout(function () {
            connected = ScatterJS.scatter.connect('hello-scatter')
            console.log(connected);
        }, 10)
    }

    //connect_scatter();
    checkLogin();
    
    //
	var selected_number = $("#selected-number").val();
	calReadyBet(selected_number,0)

    // progressbar.js@1.0.0 version is used
    // Docs: http://progressbarjs.readthedocs.org/en/1.0.0/
    var cpu = new ProgressBar.Circle("#cpu", {
        color: '#FFEA82',
        trailColor: '#eee',
        trailWidth: 1,
        duration: 1400,
        easing: 'bounce',
        strokeWidth: 6,
        from: {
            color: '#FFEA82',
            a: 0
        },
        to: {
            color: '#ED6A5A',
            a: 1
        },
        // Set default step function for all animate calls
        step: function (state, circle) {
            circle.path.setAttribute('stroke', state.color);
            var value = Math.round(circle.value() * 100);
            circle.setText("CPU: " + value + "%");
        }
    });
    var net = new ProgressBar.Circle("#net", {
        color: '#FFEA82',
        trailColor: '#eee',
        trailWidth: 1,
        duration: 1400,
        easing: 'bounce',
        strokeWidth: 6,
        from: {
            color: '#FFEA82',
            a: 0
        },
        to: {
            color: '#ED6A5A',
            a: 1
        },
        // Set default step function for all animate calls
        step: function (state, circle) {
            circle.path.setAttribute('stroke', state.color);
            var value = Math.round(circle.value() * 100);
            circle.setText("NET: " + value + "%");
        }
    });
    net.text.style.fontSize = '10px';
    cpu.text.style.fontSize = '10px';
    net.animate(0); // Number from 0.0 to 1.0
    cpu.animate(0); // Number from 0.0 to 1.0


    function countdown() { //倒计时
    	/*
        var end_time = (parseInt(Date.now() / 1000 / 86400) + 1) * 86400 - 8 * 3600; //终止时间
        var curr_time = parseInt(Date.parse(new Date()) / 1000);
        var diff_time = parseInt(end_time - curr_time); // 倒计时时间差
        var h = Math.floor(diff_time / 3600);
        h = h > 9 ? h : '0' + h;
        var m = Math.floor((diff_time / 60 % 60));
        m = m > 9 ? m : '0' + m;
        var s = Math.floor((diff_time % 60));
        s = s > 9 ? s : '0' + s;
        */
    	
    	/*
    	var open_time = new Date((parseInt(Date.now()/600000) + 1) * 600000);
    	var h = open_time.getHours()
    	h = h > 9 ? h : '0' + h;
    	var m = open_time.getMinutes()
    	m = m > 9 ? m : '0' + m;
    	var s = open_time.getSeconds()
    	s = s > 9 ? s : '0' + s;
    	*/
    	
        var end_time = (parseInt(Date.now() / 600000) + 1) * 600; //终止时间
        var curr_time = parseInt(Date.now() / 1000); 
        var diff_time = parseInt(end_time - curr_time); // 倒计时时间差
        var h = Math.floor(diff_time / 3600);
        h = h > 9 ? h : '0' + h;
        var m = Math.floor((diff_time / 60 % 60));
        m = m > 9 ? m : '0' + m;
        var s = Math.floor((diff_time % 60));
        s = s > 9 ? s : '0' + s;

        $('.timer').html(h + ":" + m + ":" + s);
        $("#bonus-remain-time").html(h + ":" + m + ":" + s)
        if (diff_time <= 0) {
            $("#bonus-remain-time").html("00:00:00")
            $('.timer').html(0 + ":" + 0 + ":" + 0);
        };


    }
    countdown();
    var start_time = setInterval(function () {
        countdown()
    }, 1000);
    
    window.onunload = onunload_handler;
    function onunload_handler() {
    	scatter.forgetIdentity()
    }
    
	Date.prototype.Format = function (fmt) { //author: meizz
		var o = {
			"M+": this.getMonth() + 1, //月份
			"d+": this.getDate(), //日
			"H+": this.getHours(), //小时
			"m+": this.getMinutes(), //分
			"s+": this.getSeconds(), //秒
			"q+": Math.floor((this.getMonth() + 3) / 3), //季度
			"S": this.getMilliseconds() //毫秒
		};
		if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
		for (var k in o)
			if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
		
		return fmt;
	}

})(jQuery);

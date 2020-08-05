/**
 * *  @alanogoddess
 * *  copyright @alano
 * */
#include <utility>
#include <vector>
#include <string>
#include <eosiolib/eosio.hpp>
#include <eosiolib/asset.hpp>
#include <eosiolib/transaction.hpp>
#include <eosiolib/crypto.h>


using namespace eosio;
using namespace std;

extern "C" {
/*
 * action parameter
 */
struct transferargs {
	account_name      from;
	account_name      to;
	asset             quantity;
	string            memo;
};
struct withdrawargs {
	account_name      user;
	asset             quantity;
};
struct openbetargs {
};
struct offerbetargs {
	account_name      user;
	asset             quantity;
	uint64_t          number;
};
struct closebetargs {
};
struct resetargs {
};
/*
 * table struct
 */
struct betaccount {
	int64_t           id;
	account_name      user;
	asset             quantity;
	account_name      inviter;
	uint64_t primary_key() const { return id; }
	account_name by_user() const { return user; }
};
typedef multi_index< N(betaccount), betaccount,
		indexed_by< N(byuser), const_mem_fun<betaccount, account_name, &betaccount::by_user> > > betaccount_index;

struct globalindex {
	uint64_t          id;
	uint64_t          gindex;
	uint64_t primary_key()const { return id; }
};
typedef multi_index< N(globalindex), globalindex> globalindex_index;

struct betstate {
	uint64_t          round;
	uint64_t          users;
	uint64_t          bets;
	asset             quantity;
	asset             inviteamount;
	asset             tpoolamount;
	asset             apoolamount;
	uint64_t          blockprefix;
	uint64_t          blocknum;
	uint64_t          result;
	uint64_t          btime;
	uint64_t          otime;
	asset             prizeamount;
	uint64_t primary_key()const { return round; }
};
typedef multi_index< N(betstate), betstate> betstate_index;

struct betoffer {
	uint64_t          id;
	uint64_t          round;
	account_name      user;
	asset             quantity;
	uint64_t          bets;
	uint64_t          number;
	uint64_t primary_key() const {return id;}
};
typedef multi_index<N(betoffer), betoffer> betoffer_index;

struct betnumber {
	uint64_t          id;
	uint64_t          number;
	uint64_t primary_key()const { return id; }
};
typedef multi_index< N(betnumber), betnumber> betnumber_index;

struct numberpool {
	uint64_t          id;
	uint64_t          bets;
	uint64_t          users;
	uint64_t primary_key() const { return id;}
};
typedef multi_index< N(numberpool), numberpool> numberpool_index;

struct innerledger {
	account_name      user;
	uint64_t          round;
	uint64_t          bounds;
	account_name primary_key() const {return user;}
};
typedef multi_index< N(innerledger), innerledger> innerledger_index;

struct innercounter {
	uint64_t          id;
	uint64_t          counter;
	uint64_t primary_key() const { return id; }
};
typedef multi_index< N(innercounter), innercounter> innercounter_index;

struct innerexist {
	uint64_t           id;
	uint128_t          exist;
	uint64_t primary_key() const { return id; }
	uint128_t by_exist() const { return exist; }
};
typedef multi_index< N(innerexist), innerexist,
		indexed_by< N(byexist), const_mem_fun<innerexist, uint128_t, &innerexist::by_exist> > > innerexist_index;

void init_betaccount(account_name code, account_name scope) {
	betaccount_index betaccounts(code, scope);
	betaccounts.emplace(scope, [&](auto& info) {
		info.id        = betaccounts.available_primary_key();
		info.user      = N(alanosystem1);
		info.quantity  = asset(0, symbol_type(S(4, EOS)));
		info.inviter   = N(none);
	});
}
void clear_betaccount(account_name code, account_name scope) {
	betaccount_index betaccounts(code, scope);
	auto betaccount_itr = betaccounts.begin();
	while(betaccount_itr != betaccounts.end()) {
		betaccount_itr = betaccounts.erase(betaccount_itr);
	}
}
void clear_globalindex(account_name code, account_name scope) {
	globalindex_index globalindexs(code, scope);
	auto globalindex_itr = globalindexs.begin();
	while(globalindex_itr != globalindexs.end()) {
		globalindex_itr = globalindexs.erase(globalindex_itr);
	}
}
void clear_betstate(account_name code, account_name scope) {
	betstate_index betstates(code, scope);
	auto betstate_itr = betstates.begin();
	while(betstate_itr != betstates.end()) {
		betstate_itr = betstates.erase(betstate_itr);
	}
}
void clear_betoffer(account_name code, account_name scope) {
	betoffer_index betoffers(code, scope);
	auto betoffer_itr = betoffers.begin();
	while(betoffer_itr != betoffers.end()) {
		betoffer_itr = betoffers.erase(betoffer_itr);
	}
}
void clear_betnumber(account_name code, account_name scope) {
	betnumber_index betnumbers(code, scope);
	auto betnumber_itr = betnumbers.begin();
	while(betnumber_itr != betnumbers.end()) {
		betnumber_itr = betnumbers.erase(betnumber_itr);
	}
}
void init_numberpool(account_name code, account_name scope) {
	numberpool_index numberpools(code, scope);
	auto numberpool_itr = numberpools.begin();
	while(numberpool_itr != numberpools.end()) {
		numberpool_itr = numberpools.erase(numberpool_itr);
	}
	for(int i = 0;i < 12;i ++) {
		numberpools.emplace(scope, [&](auto& info) {
			info.id              = i + 1;
			info.bets            = 0;
			info.users           = 0;
		});
	}
}
void clear_numberpool(account_name code, account_name scope) {
	numberpool_index numberpools(code, scope);
	auto numberpool_itr = numberpools.begin();
	while(numberpool_itr != numberpools.end()) {
		numberpools.modify(numberpool_itr, 0, [&](auto& info) {
			info.bets      = 0;
			info.users     = 0;
		});
		numberpool_itr ++;
	}
}
void destroy_numberpool(account_name code, account_name scope) {
	numberpool_index numberpools(code, scope);
	auto numberpool_itr = numberpools.begin();
	while(numberpool_itr != numberpools.end()) {
		numberpool_itr = numberpools.erase(numberpool_itr);
	}
}
void clear_innerledger(account_name code, account_name scope) {
	innerledger_index innerledgers(code, scope);
	auto innerledger_itr = innerledgers.begin();
	while(innerledger_itr != innerledgers.end()) {
		innerledger_itr = innerledgers.erase(innerledger_itr);
	}
}
void init_innercounter(account_name code, account_name scope) {
	innercounter_index innercounters(code, scope);
	auto innercounter_itr = innercounters.begin();
	while(innercounter_itr != innercounters.end()) {
		innercounter_itr = innercounters.erase(innercounter_itr);
	}
	innercounters.emplace(scope, [&](auto& info) {
		info.id    = N(ledger);
		info.counter   = 0;
	});
}
void clear_innercounter(account_name code, account_name scope) {
	innercounter_index innercounters(code, scope);
	auto innercounter_itr = innercounters.begin();
	while(innercounter_itr != innercounters.end()) {
		innercounters.modify(innercounter_itr, 0, [&](auto& info) {
			info.counter      = 0;
		});
		innercounter_itr ++;
	}
}
void destroy_innercounter(account_name code, account_name scope) {
	innercounter_index innercounters(code, scope);
	auto innercounter_itr = innercounters.begin();
	while(innercounter_itr != innercounters.end()) {
		innercounter_itr = innercounters.erase(innercounter_itr);
	}
}
void clear_innerexist(account_name code, account_name scope) {
	innerexist_index innerexists(code, scope);
	auto innerexist_itr = innerexists.begin();
	while(innerexist_itr != innerexists.end()) {
		innerexist_itr = innerexists.erase(innerexist_itr);
	}
}

// The apply method implements the dispatch of events to this contract
void apply( uint64_t receiver, uint64_t code, uint64_t action ) {
	auto _self = receiver;

	//
	if(code == N(eosio.token) && action == N(transfer)) {
		transferargs args = unpack_action_data<transferargs>();
		if(args.to != _self) {
			return;
		}

		account_name from = args.from;
		account_name to = args.to;
		asset quantity = args.quantity;
		string memo = args.memo;
		symbol_type sym = quantity.symbol;

		eosio_assert( (memo.size() <= 256), "memo has more than 256 bytes" );
		eosio_assert( quantity.is_valid(), "quantity is invalid" );
		eosio_assert( (quantity.amount > 0), "must deposit positive quantity" );
		eosio_assert( sym.is_valid(), "symbol type is invalid" );
		eosio_assert( (sym == symbol_type(S(4, EOS))), "symbol type is not EOS" );

		//
		betaccount_index betaccounts(_self, _self);
		auto betaccounts_byuser_index = betaccounts.get_index<N(byuser)>();
		auto cur_account_itr = betaccounts_byuser_index.find(from);
		if(cur_account_itr == betaccounts_byuser_index.end()) {
			//
			account_name inviter = string_to_name(memo.c_str());
			auto cur_inviter_itr = betaccounts_byuser_index.find(inviter);
			if(cur_inviter_itr == betaccounts_byuser_index.end()) {
				inviter = N(none);
			}
			//
			betaccounts.emplace(_self, [&](auto& info) {
				info.id       = betaccounts.available_primary_key();
				info.user     = from;
				info.inviter  = inviter;
				info.quantity = quantity;
			});
		} else {
			betaccounts_byuser_index.modify(cur_account_itr, 0, [&](auto& info) {
				info.quantity = info.quantity + quantity;
			});
		}
	}
	if(code == _self && action == N(withdraw)) {
		withdrawargs args = unpack_action_data<withdrawargs>();
		if(args.user == _self) {
			return;
		}

		account_name user = args.user;
		asset quantity = args.quantity;
		symbol_type sym = quantity.symbol;

		require_auth(user);
		eosio_assert( quantity.is_valid(), "quantity is invalid" );
		eosio_assert( (quantity.amount > 0), "must deposit positive quantity" );
		eosio_assert( sym.is_valid(), "symbol type is invalid" );
		eosio_assert( (sym == symbol_type(S(4, EOS))), "symbol type is not EOS" );

		//
	    asset fee(3000, symbol_type(S(4, EOS)));
		betaccount_index betaccounts(_self, _self);
		auto betaccounts_byuser_index = betaccounts.get_index<N(byuser)>();
		auto cur_account_itr = betaccounts_byuser_index.find(user);
		eosio_assert(cur_account_itr != betaccounts_byuser_index.end(), "unknow account");
		eosio_assert(cur_account_itr->quantity >= quantity, "balance is insufficient");
		eosio_assert(quantity > fee, "balance is insufficient");
		betaccounts_byuser_index.modify(cur_account_itr, 0, [&](auto& info) {
			info.quantity      = info.quantity - quantity;
		});

		auto alanosystem_account_itr = betaccounts_byuser_index.find(N(analosystem1));
		betaccounts_byuser_index.modify(alanosystem_account_itr, 0, [&](auto& info) {
			info.quantity      = info.quantity + fee;
		});

		//
		eosio::action(
			permission_level{ _self, N(active) }, N(eosio.token), N(transfer),
			std::make_tuple(_self, user, (quantity - fee), std::string("withdraw"))
		).send();
	}
	if(code == _self && action == N(openbet)) {
		require_auth(N(alanosystem2));

		//
		globalindex_index globalindexs(_self, _self);
		auto cur_globalindex_itr = globalindexs.find(N(bet));
		if(cur_globalindex_itr == globalindexs.end()) {
			cur_globalindex_itr  = globalindexs.emplace(_self, [&](auto& info) {
			info.id              = N(bet);
			info.gindex          = 0;
			});
			init_numberpool(_self, _self);
			init_innercounter(_self, _self);
			init_betaccount(_self, _self);
		}
		uint64_t bet_index = cur_globalindex_itr->gindex;

		//
		betstate_index betstates(_self, _self);
		auto cur_betstate_itr = betstates.find(bet_index);
		if(cur_betstate_itr == betstates.end()) {
			cur_betstate_itr     = betstates.emplace(_self, [&](auto& info) {
			info.round           = bet_index;
			info.users           = 0;
			info.bets            = 0;
			info.quantity        = asset(0, symbol_type(S(4, EOS)));
			info.inviteamount    = asset(0, symbol_type(S(4, EOS)));
			info.tpoolamount     = asset(0, symbol_type(S(4, EOS)));
			info.apoolamount     = asset(0, symbol_type(S(4, EOS)));
			info.blockprefix     = 0;
			info.blocknum        = 0;
			info.result          = 0x01FFFFFFFFFFFFFF;
			info.btime           = current_time();
			info.otime           = 0;
			info.prizeamount     = asset(0, symbol_type(S(4, EOS)));
			});
		}

		//
		if(cur_betstate_itr->result & 0x0100000000000000) {
			return;
		}

		// next round
		bet_index ++;
		globalindexs.modify(cur_globalindex_itr, 0, [&](auto& info) {
			info.gindex = bet_index;
		});

		betstates.emplace(_self, [&](auto& info) {
			info.round           = bet_index;
			info.users           = 0;
			info.bets            = 0;
			info.quantity        = asset(0, symbol_type(S(4, EOS)));
			info.inviteamount    = asset(0, symbol_type(S(4, EOS)));
			info.tpoolamount     = asset(0, symbol_type(S(4, EOS)));
			info.apoolamount     = cur_betstate_itr->apoolamount;
			info.blockprefix     = 0;
			info.blocknum        = 0;
			info.result          = 0x01FFFFFFFFFFFFFF;
			info.btime           = current_time();
			info.otime           = 0;
			info.prizeamount    = asset(0, symbol_type(S(4, EOS)));
		});

		clear_betoffer(_self, _self);
		clear_betnumber(_self, _self);
		clear_numberpool(_self, _self);
		clear_innerledger(_self, _self);
		clear_innercounter(_self, _self);
		clear_innerexist(_self, _self);
	}
	if(code == _self && action == N(offerbet)) {
		offerbetargs args = unpack_action_data<offerbetargs>();
		account_name user = args.user;
		asset quantity = args.quantity;
		uint64_t number = args.number;
		symbol_type sym = quantity.symbol;

		require_auth(user);
		eosio_assert( quantity.is_valid(), "quantity is invalid" );
		eosio_assert( (quantity.amount > 0), "must deposit positive quantity" );
		eosio_assert( sym.is_valid(), "symbol type is invalid" );
		eosio_assert( (sym == symbol_type(S(4, EOS))), "symbol type is not EOS" );
		eosio_assert( (number >= 1 && number <= 12), "bet number is invalid");

		//
		betaccount_index betaccounts(_self, _self);
		auto betaccounts_byuser_index = betaccounts.get_index<N(byuser)>();
		auto cur_account_itr = betaccounts_byuser_index.find(user);
		eosio_assert((cur_account_itr != betaccounts_byuser_index.end()), "unknow account");
		eosio_assert((cur_account_itr->quantity >= quantity), "balance is insufficient");
		betaccounts_byuser_index.modify(cur_account_itr, 0, [&](auto& info) {
			info.quantity      = info.quantity - quantity;
		});

		//
		globalindex_index globalindexs(_self, _self);
		auto cur_globalindex_itr = globalindexs.find(N(bet));
		eosio_assert((cur_globalindex_itr != globalindexs.end()), "global index is not exist");
		uint64_t bet_index = cur_globalindex_itr->gindex;

		//
		betstate_index betstates(_self, _self);
		auto cur_betstate_itr = betstates.find(bet_index);
		eosio_assert((cur_betstate_itr != betstates.end()), "bet state is not exist");
		eosio_assert((cur_betstate_itr->result & 0x0100000000000000) >> 32, "current bet is not opening");

		//
		numberpool_index numberpools(_self, _self);
		auto cur_numberpool_itr = numberpools.find(number);
		eosio_assert((cur_numberpool_itr != numberpools.end()), "system is not initilize");
		uint64_t bets = cur_numberpool_itr->bets;
		uint64_t users = cur_numberpool_itr->users;

		//
		uint128_t numberuser = number;
		numberuser = (numberuser << 64 ) | uint64_t(user);
		innerexist_index innerexists(_self, _self);
		auto innerexist_exist_index = innerexists.get_index<N(byexist)>();
		auto innerexist_itr = innerexist_exist_index.find(numberuser);
		bool numberuser_exist = (innerexist_itr != innerexist_exist_index.end());
		if(!numberuser_exist) {
			if(number == 12 || number == 11) {
				eosio_assert(users == 0, "users of number 11 and 12 is limited to 1");
			} else if(number == 10 || number == 9) {
				eosio_assert(users < 2, "users of number 9 and 10 is limited to 2");
			} else if(number == 8 || number == 7) {
				eosio_assert(users < 4, "users of number 7 and 8 is limited to 4");
			} else if(number == 6 || number == 5) {
				eosio_assert(users < 8, "users of number 5 and 6 is limited to 8");
			} else if(number == 4 || number == 3) {
				eosio_assert(users < 16, "users of number 3 and 4 is limited to 16");
			} else if(number == 2 || number == 1) {
				eosio_assert(users < 32, "users of number 1 and 2 is limited to 32");
			}
			innerexists.emplace(_self, [&](auto& info){
				info.id           = innerexists.available_primary_key();
				info.exist        = numberuser;
			});
		}

		//
		asset quantity_prebet(10000, symbol_type(S(4, EOS)));
		uint64_t this_bets = 0;
		asset this_quantity(0, symbol_type(S(4, EOS)));
		betnumber_index betnumbers(_self, _self);
		while(quantity >= quantity_prebet) {
			quantity = quantity - quantity_prebet;
			betnumbers.emplace(_self, [&](auto& info){
				info.id             = betnumbers.available_primary_key();
				info.number         = number;
			});
			this_quantity += quantity_prebet;
			this_bets ++;
		}

		//
		betoffer_index betoffers(_self, _self);
		betoffers.emplace(_self, [&](auto& info){
			info.id           = betoffers.available_primary_key();
			info.round        = bet_index;
			info.user         = user;
			info.quantity     = this_quantity;
			info.bets         = this_bets;
			info.number       = number;
		});

		//
		betstates.modify( cur_betstate_itr, 0, [&](auto& info) {
			info.bets          = info.bets + this_bets;
			info.quantity      = info.quantity + this_quantity;
			info.users         = info.users + (numberuser_exist ? 0 : 1);
		});

		//
		numberpools.modify( cur_numberpool_itr, 0, [&](auto& info) {
			info.bets          = info.bets + this_bets;
			info.users         = info.users + (numberuser_exist ? 0 : 1);
		});

		//
		account_name inviter = cur_account_itr->inviter;
		if(inviter == N(none)) {
			return;
		}

		//
		innerledger_index innerledgers(_self, _self);
		auto invite_innerledger_itr = innerledgers.find(inviter);
		if(invite_innerledger_itr == innerledgers.end()) {
			innerledgers.emplace(_self, [&](auto& info) {
				info.user       = inviter;
				info.round      = bet_index;
				info.bounds     = this_bets;
			});
		} else {
			innerledgers.modify( invite_innerledger_itr, 0, [&](auto& info) {
				info.bounds     = info.bounds + this_bets;
			});
		}

		//
		innercounter_index innercounters(_self, _self);
		auto innercounter_itr = innercounters.find(N(ledger));
		innercounters.modify( innercounter_itr, 0, [&](auto& info) {
			info.counter         = info.counter + this_bets;
		});
	}
	if(code == _self && action == N(closebet)) {
		require_auth(N(alanosystem2));
		closebetargs args = unpack_action_data<closebetargs>();
		//
		globalindex_index globalindexs(_self, _self);
		auto cur_globalindex_itr = globalindexs.find(N(bet));
		eosio_assert((cur_globalindex_itr != globalindexs.end()), "global index is not exist");
		uint64_t bet_index = cur_globalindex_itr->gindex;

		//
		betstate_index betstates(_self, _self);
		auto cur_betstate_itr = betstates.find(bet_index);
		eosio_assert((cur_betstate_itr != betstates.end()), "bet state is not exist");
		eosio_assert((cur_betstate_itr->result & 0x0100000000000000) >> 32, "current bet is not opening");
		if(cur_betstate_itr->bets == 0) {
			return;
		}

		//
		uint64_t otime = current_time();
		uint64_t btime = cur_betstate_itr->btime;
		uint64_t mtime = ((btime / 600000000) + 1) * 600000000;
		eosio_assert((otime >= mtime && otime % 600000000 <= 10000000), "time windows of reveal is invalid");

		//
		checksum256 hash;
		int block_prefix = tapos_block_prefix();
		int block_num = tapos_block_num();
		auto mixed_block = block_prefix * block_num;
		const char* mixed_char = reinterpret_cast<const char*>(&mixed_block);
		sha256((char*)(mixed_char), sizeof(mixed_char), &hash);
		const char* p64 = reinterpret_cast<const char*>(&hash);
		uint64_t bet_result = (abs((int64_t)p64[5]));
		bet_result = bet_result % cur_betstate_itr->bets;
		/*
		uint64_t bet_result = hash.hash[1];
		bet_result = bet_result << 32;
		bet_result |= hash.hash[0];
		bet_result = bet_result % cur_betstate_itr->bets;
		*/
		//uint64_t bet_result = N(hash.hash);
		//bet_result = bet_result % cur_betstate_itr->total;
		//uint64_t bet_result = (hash.hash[15]) % cur_betstate_itr->total;

		// prize
		asset bet_quantity = cur_betstate_itr->quantity;
		asset bet_apoolamount = cur_betstate_itr->apoolamount;
		asset quantity_prebet(10000, symbol_type(S(4, EOS)));
		betnumber_index betnumbers(_self, _self);
		auto prize_betnumber_itr = betnumbers.find(bet_result);
		uint64_t prize_number = prize_betnumber_itr->number;
		numberpool_index numberpools(_self, _self);
		auto prize_numberpool_itr = numberpools.find(prize_number);
		uint64_t prize_bets = prize_numberpool_itr->bets;
		asset prize_total = bet_quantity + (bet_apoolamount / 2);
		asset prize_gains = bet_quantity - quantity_prebet * (prize_bets);
		innercounter_index innercounters(_self, _self);
		auto ledger_counter_itr = innercounters.find(N(ledger));
		uint64_t all_bounds = ledger_counter_itr->counter;
		asset bet_inviteamount = all_bounds == 0 ? asset(0, symbol_type(S(4, EOS))) : (prize_gains / 10) * 2;
		asset bet_tpoolamount = (prize_gains / 10) * 1;
		bet_apoolamount = bet_tpoolamount + (bet_apoolamount / 2);
		asset prize_amount = prize_total - bet_inviteamount - bet_tpoolamount;
		asset prize_amount_prebet = prize_amount / prize_bets;

		//
		betaccount_index betaccounts(_self, _self);
		auto betaccounts_byuser_index = betaccounts.get_index<N(byuser)>();

		if(bet_inviteamount.amount > 0) {
			innerledger_index innerledgers(_self, _self);
			asset invite_amount_prebound = bet_inviteamount / all_bounds;
			auto invite_innerledger_itr = innerledgers.begin();
			while(invite_innerledger_itr != innerledgers.end()) {
				auto cur_account_itr = betaccounts_byuser_index.find(invite_innerledger_itr->user);
				betaccounts_byuser_index.modify(cur_account_itr, 0, [&](auto& info) {
					info.quantity     = info.quantity + (invite_amount_prebound * invite_innerledger_itr->bounds);
				});
				invite_innerledger_itr ++;
			}
		}

		//
		betoffer_index betoffers(_self, _self);
		auto betoffer_itr = betoffers.begin();
		while(betoffer_itr != betoffers.end()) {
			if(betoffer_itr->number == prize_number) {
				auto cur_account_itr = betaccounts_byuser_index.find(betoffer_itr->user);
				betaccounts_byuser_index.modify(cur_account_itr, 0, [&](auto& info) {
					info.quantity     = info.quantity + (prize_amount_prebet * betoffer_itr->bets);
				});
			}
			betoffer_itr ++;
		}

		//
		betstates.modify( cur_betstate_itr, 0, [&](auto& info) {
			info.blockprefix    = block_prefix;
			info.blocknum       = block_num;
			info.result         = prize_number;
			info.otime          = current_time();
			info.prizeamount    = prize_amount;
			info.apoolamount    = bet_apoolamount;
			info.inviteamount   = bet_inviteamount;
			info.tpoolamount    = bet_tpoolamount;
		});
	}
	if(code == _self && action == N(reset)) {
		require_auth(_self);

		clear_betaccount(_self, _self);
		clear_globalindex(_self, _self);
		clear_betstate(_self, _self);
		clear_betoffer(_self, _self);
		clear_betnumber(_self, _self);
		destroy_numberpool(_self, _self);
		clear_innerledger(_self, _self);
		destroy_innercounter(_self, _self);
		clear_innerexist(_self, _self);
	}
}
} // extern "C"

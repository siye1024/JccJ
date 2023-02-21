package usersrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/db"
	"dousheng/pkg/jwt"
	user "dousheng/rpcserver/kitex_gen/user"
	svr "dousheng/rpcserver/kitex_gen/user/usersrv"
	api2 "dousheng/rpcserver/user/api"
	"errors"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/gorm"
	"log"
	"net"
	"time"
)

// Arg2Config Directly set the parameters
var (
	Arg2Config = &api2.Argon2Params{
		Memory:      65536,
		Iterations:  3,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Register Success"
	)

	// empty username or password has been processed by user client
	if len(req.Username) == 0 || len(req.Password) == 0 {
		err := kerrors.NewBizStatusError(10001, "Empty Username or Password")
		return nil, err
	}
	// len of username or password exceed the 32 bit
	if len(req.Username) > 32 || len(req.Password) > 32 {
		err := kerrors.NewBizStatusError(10014, "Too Long Username or Password")
		return nil, err
	}

	err = api2.NewCreateUserOp(ctx).CreateUser(req, Arg2Config)
	if err != nil {
		return nil, err
	}

	//Auto Login
	uid, err := api2.NewCheckUserOp(ctx).CheckUser(req)
	if err != nil {
		return nil, err
	}

	token, err := xhttp.Jwt.CreateToken(jwt.CustomClaims{ //Claim is payload
		Id:   int64(uid),
		Time: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	resp = &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserId:     uid,
		Token:      token, // successful resp must have token
	}

	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Login Success"
		respErrorMag  = "Signature is Invalid"
	)
	// empty username or password has been processed by dousheng client
	if len(req.Username) == 0 || len(req.Password) == 0 {
		err := kerrors.NewBizStatusError(10001, "Empty Username or Password")
		return nil, err
	}
	// len of username or password exceed 32 bit
	if len(req.Username) > 32 || len(req.Password) > 32 {
		err := kerrors.NewBizStatusError(10014, "Too Long Username or Password")
		return nil, err
	}

	// check the user's information
	uid, err := api2.NewCheckUserOp(ctx).CheckUser(req)
	if err != nil {
		resp = &user.DouyinUserRegisterResponse{
			StatusCode: 10012,
			StatusMsg:  &respErrorMag,
		}
		return resp, nil
	}

	token, err := xhttp.Jwt.CreateToken(jwt.CustomClaims{ //Claim is payload
		Id:   int64(uid),
		Time: time.Now().Unix(),
	})
	if err != nil {
		resp = &user.DouyinUserRegisterResponse{
			StatusCode: 10012,
			StatusMsg:  &respErrorMag,
		}
		return resp, nil
	}

	resp = &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserId:     uid,
		Token:      token, // successful resp must have token
	}

	return resp, nil
}

// GetUserById implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserById(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	var (
		respStatusMsg    = "Get User's Info By ID Successfully"
		avaterurl        = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAoHCBIUFBUSFhYYGRgZGBofGRwWHBwZHBYaGhwaGhwfHhgcJC4lHB8tIRocJ0YmLS8xNTU1GiQ7QDs1Py40NTEBDAwMEA8QHxISHDQsJSw/NDQxNDY0NDQ2NDQ0MTQ0NDQ0NDQ2NDQ0NDQ0NDQ0NDQxNDQ0NDQ0NDQxMTQ0MT80NP/AABEIAOEA4QMBIgACEQEDEQH/xAAcAAEAAwADAQEAAAAAAAAAAAAABQYHAQMEAgj/xABGEAACAQIDBAcEBwYEBAcAAAABAgADEQQSIQUGMUEHIlFhcYGREzJCoVJicoKSscEUI6KywtEkU9LhFUNzkxYXM1Sz8PH/xAAYAQEAAwEAAAAAAAAAAAAAAAAAAQIDBP/EACIRAQEBAAEEAgMBAQAAAAAAAAABAhEDEiExQVEEEzIiYf/aAAwDAQACEQMRAD8A2aIiAiIgIiICIiAiIgIiICIiAiIgIiICIkZt/aYwuHq4kqWyLfKNMxJCgX5akawJOJVNyN6jj1qZqYR6ZFwpJUq17EX1HAy1wESE2rvNg8NUSjWqhHcAgWYgAmwLECyi/M24GTcBERAREQEREBERAREQEREBERAREQEREBERAREQEREDiYXvBsra/wC+q4gVWRSS7FwUtm0Krmtl4WAGnZN1lR6R8bTTA16ZdQ7hAqkjM3XW9l4nS5kwZ/0b4Ks2JbEI1koKWqC5BcMjhUtwNyL68Lds82I3v2ninslRxfUJhwVsPu9Y+JMm+ivEIq45CesaasB2qgcMfIsvrKJs7B+2YU89NOre9ZwiaW0zHS8lDt2ocRnP7SXL2F/ali+W2l82tpZdib6Y7B1Vp4hnemCA61gc6KfiVmGa4GtjcEesq+0MC2Hc02KEgA3psHUgi4sy6Gevb2Bag1Om9RHcUxmyNnWndmsmbtA1tyzQP0OrAgEagi47xPueDYd/2bD34+xp38ci3nvlUkREBERAREQEREBERARE6qtVVUsxCqBcljYAd5PCB2RKNtvpIwlElKINdxzU5UB+2eP3QR3yl47pD2jWOVGWmDwWkmZvxNc+gEngbbPNXx9FPfqIv23VfzMwxsJtTEav+0OD/mOyr6OQPlPqluTjDqRTX7T/AOkGTMavwr3Rs7bw4IccTQ/7if3nA3hwR4Ymh/3E/vMhXcTE83ojzc/0TltxMTyel6uP6Zb9ejujaaGOpP7lRH+yyt+RnomB1tzMYmoRH+w4v/FaKW0Nq4Pg+IQDk13QeTZlkXFhNSt9iZBszpQxS2FamlUdqXpt+qn0EvWwd88Hi7Kr5HPwVLKx+yb2byMrwsmto06jUai02yOyMEb6LEEKfI2mG7V3Sx1ClUxWIUKAyglnDs5ZgoIte/HiSJvszLpQ3jotTbAoWNRXQvp1VABa1/pajhERXx0RYFCmJqsoLdWnc/QIzMPPq+gnO1uiwE3w1bKNepWFwOyzqL28QT3yS6IqOXB1H+nXa3gqoPzvPVvfvzTwT+wWmalTKGNzlRQb2ubEk6XsB2ax8jJdvbHq4OqaFXKWChuoSQQ17akDsMn93N1cO+ISlXxVAnQ+ypFnL6Xyl7BVPaASePDjILb+26mNrGvUCBsoUBBYBQSQNSSeJ1ly3C3QUth8bVrJr16VNT1mKk2LE9hHugHhx5SRq4FtBPqIlUkREBERAREQEREBETOekPfI0ScJh2tU/wCY440wR7qnk55nkO86BL71b74fB3pp+9rfRU6J9tuX2Rr4cZmOMx+0NpvZizgH3V6tJPLh5kkz1bu7ptXtWrXVDqF+OpfW5Pwqe3ifnL/hsPTpoKaKFUcAosJtjpc+az1vj0qeytx6a2au+c/QS6r5txPlaWrCYKlSGWmiIPqgC/ieJnfE3zmT0yurSIiWQREQERECMx+wMLXvnprmPxJ1G9Rx87ym7Z3Lq07vRPtEGuXg48uDeWvdNFiV1jOlpqxSN1N/62HIpYgtUpcLnWpT5cTqwHYdew8padtbl4XaLjGUa+X2gBYoA6vYAXAuCrWFj4cL3kZvJuzTxANRLJV7eAfubv8ArSs7pbx1dnVzTqBvZFrVUPFDwzKPpDu94eRnNvFy1zqVs2xdl08LQTD075UB1PEkklie8kkzybW3XweKdatakGZRa+ZluBqAcpGYanj2yVo1VdVdSCrAFSNQQRcEHstO6ZrqTvnuclfDKmGpoj0mLIqhUDBhZ1vwubKbnmovKTu3uTj3xFN3RqKU3Vi7EAjK2ayAG5J7eGvHt2yJPIRESAiIgIiICIiAiJwTAr++O3Rg8M9UWzt1aYPNzfW3YACfLvmUbpbHOJqtWq3ZFa7FtfaOTex7eNz4jtn30l7z0cTiVSnUV6dJbKV6ys7WLEEaHgo+6ZNbO2/hMLhKTWrezI0qewqBHc6khioBub2/2mmJOfKurePC2RIXYG82Hx1b2FAVCwUszMuVUUaXJJvqSBa3Pxnm3k3upYGucPUpVWYKrArkysrXsQSwPEEajlOjvz9se3X0scSiHpNw3+RW/g/1Q/Sbh/hoVj4lB/UY/Zn7OzX0vcSgf+Z1H/21T8aT7HSbh+dCt5FD/VH7M/ae2/S+RKQvSZg/ipVx32Qgfxyex28dLDolTEJXoo/uF6bam17HJfKba2NjoeyO/P2jt19JmJEbK3lweJf2dFy72Jyqj3sOJIK6DUa94nadvYQMUOIpq6khkdwjKRxBVrEGT3Z+0cX6SUTpo4qm/uOjfZYN+RnfJQ4lS342JnT9pQddB1wPiQc/Ffy8JbYIvpI1nunCZeLyz7dTfmvgwtJx7SgD7vB6YJucjcx9U+omwbI2tQxVMVaLhlPHkVPYy8Qe6YTvPsk4auygdR+sngeK+R08LTr2PtTE4N1r0mK5uR1WoFOqsOdvUX75yazxW8r9FRIDdfeSljqWdOq62FRCblCfzU2Nj+ukn5VYiIgIiICIiAiIgJG7xEjCYkjiKFW3jka0kpQd/d80w6thaNnqspDk6rSVhqCObkHhy4nsIY9uXu+cdiadDXIBmqsPhpra4vyLaKPG/Kad0ubXp0MGuCULmq5bLyp0qZVrgctVVR59koe5W9B2VUrq9E1A6oCL5GXLmKkZh7pD39DO7Ydb/i210fElQrEsEv1ctMFkpLfj2nt63bJGhdG+xVwOCbE1rI1VfaVC2ns6aglVN+FluxHax7Jj2822XxuJq4lrgO1kU/Ci6Kvpqe8maX0w7xZKa7PRus9mrW+FAbqp7MxF/Be+Y/AS5dH27tHFGvUroWRMqrZmXrG5b3SL2FvxSmmbFunVwmEwtOk1elnPXezqeu2pGh5Cy/dl+nJdeVNXieHd/wCBdm/5H8dT/VOttwdmm9qTDwqVP9UkX3mwS8a6eQZvyE8z744IcHY+CP8AqBN+Mf8AGf8AtROjLd/9pxoLi6YezuD8TA2RbfaBP3Lc5ZemfboPssAtiQRVqHs0KovndmPgvbKpU3gqYHGV6+CayVtSHW6m5zEFT2MWIII0adG7uzqu1MeBVYvnYvXbh+7W2bha1+qgtwuOyctnltGl9FGwxhsI2LqdV64D3OmSioJS/Ze5fwI7Jku9G01xWMxGJUWV3JT7CgIpPeVUHzmr9Le3hh8KuEQ2euLEDTJRWwbwzaL4ZuyYnCQKAb217ZYN09s4mnicOi1XKPURGRnZkKuwU9UmwNje/dK/LP0e7NatjUe3Uo9djyvqEHiTr90y2ee6cK69eWyxETsc6F3t2YMRh3sOvTBdO3Qar5j52lM3SFOsXwdTVKgzIeaVFHFew5fXLNNmUVP8JjieAp1r/cJv/IZjucWVfN5ljspVMTszFBlPWTxC1aZ5HuNvIjum47E2rTxVBMRTPVYcDxVhoVPeDpKLvXscYiiSou6AshHPtXzHzAkH0Y7f9hXOGc2p1iAL8Eq8F8Mw6vjlmW88Vpm8xs8REzXIiICIiAiIgVXf3eE4LDXQj2tQlU+rp1mt3D5kTPdy93v2h/2mtdlzEqG1zsDdma/EX9Tfs1++lHFtVx4ojhTREA+s/WPqGT0l42XSWiiUxoqIF9Bxmmc3i2K8zukvpB76bk08aoqUyErqLAn3ag5K9uHcw4X5zHtp7NxGFfLWR6bA9Um4BI5o40PiDP0iNZ2Lg/aDKygqeOYXB8jxmMt58urfTzxzLw/LzuWJYksTxLEkk95OpnAF9JvHSJuzglwVatTw1FKitTOZEVTYuqnVR2GUbd/YlPEYeoq5UqJUurW5FV0a2uU2PhLWyTmscYurxn28WG3IqU8KuPr1qSIy3ppcuzE+6AF0zHsvpz4aR1EKWUOxVSRmYDMVF9Tl0vYcpP0dz8UzBWKKo+LNmsOdlH+0ncZujSaitNDldODnXOTxz25flK3qZny1n4vUst4Ru8O62FwlBawxftWqAGkqKoDg8WvmNlA59thKqoAKlg2UnW2hYA9bKxBF++xtLJhNy6xfruipfUoSzHwBA9TLRjNh0alFaFsoQdQjih7e+/Pti9TMWx+J1NS2zhXdubP2ImFDUnq1arr+7BcZqbcLuFAC2PI8fDWVPdvb1fZlc1UVHzIUZXvldbgixGqkEfnpLPid2Ew6NVrVQyLwVBlLnktydL9152dHmwVxeKL1FDUqYzOpF1LNcItjxHE/d75aal9Md9O48VR94Ns1sbXbE1SMzAABbhUVfdVQdbanzJkbN+2n0VbMqksi1KJP+U2n4XDADwtIyl0PYRWu1as47CVW/iQt/QiOVJOWO7PwNXEVFo0kLux0A+ZJ5KOZPCbdu5sJMFRFEEM/vVHHxOePkOA7hLDsfdehhUKUESncakAszHlmdtWnmqqVYq3EHWbdHzbVOrxJJHxEROlgTOukDD5cQlTk9MX7ypIPyyzRZS+kan1cO/Yzj1Cn9Jn1J4Wx/SybAxGfDUH5lFB8VGU/MSi767K9hXFVNEqEsLaZXGrW7O3zPZLVuNUvg0HYzj+In9Z2b44UVMJU7Us692Xj/CWjU7sJl40tW5W2v2zCU6rHrr1an21tr5izfelhmWdDuM62JoHgQjjxF0b+manOStyIiAiIgIiIGF7wnPth+/E0x+H2a/pNIo0yzBRxJmb7z/utr1GOgGIpv5EU3v8AOaxsWlcs/YLDxP8A9+c2xrtxay1OdSJelSCqFHATtiJi1eHauCWvRqUG92ojKe64tfxHHymHYDE1dnYl6dVDoctRRzA1Vlvx7R2gzfjIPeHdnC41bVUswFlddHXuvzHcbiPF8VOdXOpqe4q2B2pQrC6VFPdezDxU6z2ytbQ6LcSpJo1abjkHvTb5BgflI47kbXXQI1vq1kA+biZXpT4ruz+dePMXOpVVRdmCjtYgfnITH704ancK3tG7E4eb8PS8jMP0cbScgv7NO0u+Yj8Ia/rLJsrouorY4iq1T6qD2a+ZuWPkRJnSk91XX5ur/M4UhFxm1K4pot7chcJSB+Jm/XieQ5TZN29h08FQWgmp4u1rF3PFj6AAcgBPZs/Z1Gggp0aa01HJRa57SeJPedZ65p/yOPWrq82+X1ERCHEjdq4XMucDUfMSTnEnOrLzEWczhUYnox1HI7KOHEeBnnnbLzOXPZwSndIx/d0B9dv5ZcZS+kdurhx31D6Bf7yu/wCanPtI7hj/AAg+2/5iSm3MRTTD1Wc9Uoy95LAgAd5JkRulXSlgVqOwVQ1Qkn7RHmdOEqu2tq1cdWSmiMVzWpoOLE8zyv8AIDzMrdSZie3nSe6IcOxxdV9bLRIPYS7Jb+Q+k2GV3c7d5cFhwmhqN1qjDm3YPqqNB5nnLFOWtyIiAiIgIiIGWdLWxDmTGqNLBKluXHIx9Sv4ZPdG23VxFA0mb99T98Hiy2ADjt4WPf4iW3GYVKqNScBkdSGB5gzFNv7CxWysQtakzZM37uqo4fUflflY6N6gTL44Q3OJQN2OkalXZKNdfZ1GIUMuqOx0HehJ7bjvl/kJIiICIkLvFvDh8EgeqTdr5EXVnI7B2ajU6C8CZkRtTeXBYYkVa9NWHwg5n/At2+UyTbu/GNxZyITSQmwSkTmbuZh1mPcLDunn2dubiqnWcLTU/T1b8A/UiWmLVbqRoGJ6UMEpsiVn7wqqP4mB+U869KuF50K/8B/qkJQ3DogderUY/VCoPQhvznc+42EPB6o+8p/NZf8AVUd8WbCdI+zn0ZqlP/qIbeqZpY9n7Vw1cXo1UqduRgxHiBqPOZVidwh8FY+DqD/EpH5SCxu7eMw59oFJy6h6JJK9+lmHjaVvT1PhM1K/QMTGd3OkbE0SqYj99T+loKijubg/gde+atsnatDE0xVouHU8bcVPYVOqnuMpws8m3F6yntB+R/3kXJTbjAsq9gJ9f/yRc6+n/Mc+/wCiU3pGp9Sg3Y7D1UH+mXKVrf2mDhb/AEXQjzuv6ydz/NM+1EwiYnEGnhKYZ7E5EXgLm7MeQ4+8eE2Lc3c+ngV9o1nrsLM3JQfhS/Adp4n5SA6HEXLimt1syC/dZjb1vNLnJa3kcxESEkREBERAREQEpnSftNaOCanYFq5CAHkB1mbysPMiXOYv0l41sRj1w6a+zCoo+u5DMfmg+7JiK53S3apvTp4mpmzZyyKDZbKeqTpc6gnj2T0b67Q2kaq06bVvZlRl9lmGZtc2Zk18iZasLh1polNeCKFHkLTunV+udvDHuvPLK/8Ahm0veyYjxzNf87ztobw7TwpA9tXX6ta7g+VQH5TT581EVgVYBgeIYAg+Rlb0Z9p/ZVb2V0puLLiKIb61E2P4GNj+ISob0baOMxTViSEvlQN8CDhcC+p1Y25mXTHbn4Opcqppn6hsPwG49LSHbcE3Nq4y8rpr59aUvSsW74+ti7Y2ZhlATOz26zlDdjztfVR3CWjZ+2cNiNKdRWP0TdW/C1jK2Nwkt/67X+wLel5A7Y3ZxGG/eDroNc6XBTvI4r4i475eXWZ68K8TXy1GJQ9hb6MlqeJuy8nAuw+0Pi8Rr4y8YeujqHRgyngVNwZfOpr0rrNjsiIl1UDt3dijiLutkqfTUaMfrrz8eMo9OpjNm18yko/aNUqKOR5OvzF+RmrTrr0EcZXRXF72cBhfwMz105peasdez9oNiKdPEMuVnUMRyHLS/LT0nogCJeTicK0lf35H+Df7dP8AmEsEh97qWfB1h2AN+FlJ+QMjf80z7jo6Gz1cWO+l8xU/tNOmTdDte1bE0/pIjfgZgf8A5BNZnHXSRESAiIgIiICIiB8E2BJ5X9Jhm7THE7Ras2t2qVNe8nL6Zh6TYt4sR7PCYmoOK0ahHiFa3zmTdHVEF67diKv4ib/yy/TnOorq+F3oY2k7MiOjMvvKrAlfEDhO+ZbtXYeLwLrVGYJfqVU4cbWYj3W7jx75YNib6IwFPE9VvpqOqftKPdPeNPCb56k54rK5+YuUTqw2Jp1BmR1cdqEMPlO600UcROnE4qnTF3dEHa7Bfzld2lvph0uKYNVu7qp+I6nyEi6k9pkt9LOzAAkmwHEngJSt5t7RZqOHa99GqDgBzCdv2vTtkBjdrYvGsKYzNmPVpUgbHxAuW8T8pdN1ujexWrjLdooqbj77Dj9kadpPCY76vPiNM4+1Z3c3IxOMovWUhFA/d5wbVW52I4L9ax1854KOJxez6zIysjD3kcXVx29hH1lPnP0DTQKAAAABYAaAAcAByE8W19j4fFJ7OsgccidGU9qsNVPhMprhpwoOxt68PXsrH2b/AEXPVJ+q/A+BsZYJTt4OjSvTJfDN7VPoMQrjz0V/l4SsjHY/BnITVpW+GopA8g4tbwmuer9s7j6avEzNN9cYB7yHvKf2IEntz9v4jEVHSpZlC5syrlym4ABtob69+k0nUzbwpc2LdERNFSdWLoCoj0zwdGX8QI/WdsSBnHR9izQ2hSDaZy1N/FtAPxqs3Wfn7eai2HxrumhzrUQ9hJz/AMwPpN4wGKWtSp1V4VEVh4MAf1nHqcXh05vMeqIiVSREQEREBERAgN9zbZ+K/wCk3z0mddHAGXEHnen+TzUd4cIa2FxFIcXpOo8Spt87TIujrFAValM/GgI8UP8AZj6TTpf1FN+mu7Ppq9HIyhlOYFWAIIvwIPHjKVt/oypuS+FcUyfge7J5MOso7rN5S47Eq+8nmPyP6SWldzjVTnzGEYrcjaVJrigzfWpMrX9CG+U6f+BbVOnscV55rfM2m/RI7k8MJwm4W06huaOT61R1HyBLfKWnZPRaoIbE1s31KQyjzc6keAHjNNiOThHbK2Ph8MuShTVBzIHWb7THVvMyRiJCSIiAnwyAixAI7DqJ9xAjn2LhCbnD0Se000J9bSIxtJEdlRVVRYWUBRe3YJZ5WseP3jfamvR/pn1PTzRETqYkRECk9IuD0pVx3o38y/1esuHRbtD2uCCE3NF2T7ps6/JrfdkLvyB+xvfkyW8cwH5Ex0NubYscr0j5n2g/QTl6s/03xfDT4iJkuREQEREBERA4ImE7zbOfZuOzU1ypmz0TyZT7yeAuVtxtbtE3eQO9W71PG0DSayuutN7XyN/pPAj+wky8IqJ2JtVKipiEN1PEcx9JT3iXBHDAEagzMdzt1doYepU9ooVCLWDqwdgRZlAOgtfU2Oo0mjbPolEAPG5Ph3TTdmpL8q5ll4eyIiZLkREBERAREQEREDiV/a6WqE9oB/T9JYJHbVwpYBgLleXaJfp640pucxBRO1MO5Ngrek7hs6r9D5j+86rrM+WPF+nkie5dl1TyA8T/AGnemx25uPIXkXqZnynt19M+6QcYFopR+J3zHuVP9yPQyb6I9nMmGqV2/wCa4C/Zp5lv+It6T53q3BqYqslSnWVRlCsHBJWxJuoXjx4G3DjLts3BpRpU6Ke7TUKO2wHE954+c5t67ry1zniPZERKLkREBERAREQEREBERAREQEREBERAREQEREBERAREQEROCIHMREBERAREQEREBERAREQEREBERAREQEREBERAREQEREBERAREQEREBERAREQERED/2Q=="
		backgroundurl    = "https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/d2a036592c1544e6a73a3fb9508963f5~tplv-k3u1fbpfcp-watermark.image?"
		defaultsignature = "这个人在和喜羊羊嬉戏，什么都没有留下~"
	)

	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if req.UserId < 0 {
		return nil, kerrors.NewBizStatusError(10007, "Invalid Username")
	} else if req.UserId == 0 {
		req.UserId = claim.Id
	}

	u := new(db.User)
	if err := db.DB.WithContext(ctx).Where("id = ?", req.UserId).First(u).Error; err != nil {
		return nil, err
	}
	if u == nil || u.ID == 0 {
		err := kerrors.NewBizStatusError(10009, "User Already Withdraw")
		return nil, err
	}

	// true means the claim.id has follow the modelUser.id, false means not follow
	isFollow := false
	if claim.Id == int64(u.ID) {
		isFollow = true
	} else {
		relation := new(db.Relation)
		err := db.DB.WithContext(ctx).First(&relation, "user_id = ? and to_user_id = ?", claim.Id, u.ID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if relation != nil && relation.UserID != 0 { //double check is necessary
			isFollow = true
		}

	}
	works, err := db.PublishList(ctx, int64(u.ID))
	if err != nil {
		return nil, err
	}
	workCount := int64(len(works))

	favs, err := db.FavoriteList(ctx, int64(u.ID))
	if err != nil {
		return nil, err
	}
	favCount := int64(len(favs))

	userInfo := &user.User{
		Id:              int64(u.ID),
		Name:            u.UserName,
		FollowCount:     &u.FollowCount,
		FollowerCount:   &u.FollowerCount,
		IsFollow:        isFollow,
		WorkCount:       &workCount,
		FavoriteCount:   &favCount,
		TotalFavorited:  &u.FavoritedCount,
		Avatar:          &avaterurl,
		BackgroundImage: &backgroundurl,
		Signature:       &defaultsignature,
	}

	resp = &user.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		User:       userInfo,
	}
	return resp, nil
}

func (s *UserSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := svr.NewServer(new(UserSrvImpl),
		server.WithServiceAddr(addr),                            // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler), //support kerrors
		//server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "userRegisterLoginGetInfo"}),
	)

	log.Println("Start User RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *UserSrvImpl) Stop() {
	log.Println("Stop User RPC service...")
}

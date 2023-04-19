package goods

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"net/http"
	"strconv"
	"strings"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	fmt.Println(err)
	errs, ok := err.(validator.ValidationErrors)
	fmt.Println(ok)

	for _, err := range errs {
		fmt.Println(err)
	}

	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "用户已存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func List(g *gin.Context) {
	//商品的列表
	request := &proto.GoodsFilterRequest{}

	priceMin := g.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := g.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := g.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := g.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	isTab := g.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := g.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	perNums := g.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := g.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := g.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//请求商品的service服务
	rsp, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.S().Errorf("[List] 查询 [商品列表] 失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, g)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, goods := range rsp.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          goods.Id,
			"name":        goods.Name,
			"goods_brief": goods.GoodsBrief,
			"desc":        goods.GoodsDesc,
			"ship_free":   goods.ShipFree,
			"images":      goods.Images,
			"desc_images": goods.DescImages,
			"front_image": goods.GoodsFrontImage,
			"shop_price":  goods.ShopPrice,
			"category": map[string]interface{}{
				"id":   goods.Category.Id,
				"name": goods.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   goods.Brand.Id,
				"name": goods.Brand.Name,
				"logo": goods.Brand.Logo,
			},
			"is_hot":  goods.IsHot,
			"is_new":  goods.IsNew,
			"on_sale": goods.OnSale,
		})
	}
	reMap["data"] = goodsList

	g.JSON(http.StatusOK, reMap)
}

//go:build unit

package vizmedia_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/vizmedia"
	"github.com/stretchr/testify/assert"
)

type ParseVolumeHtmlTestCase struct {
	InputHtml            string
	InputSeriesName      string
	ExpectedName         string
	ExpectedRedirectLink string
	ExpectedIsReleased   bool
	ExpectError          bool
}

var (
	twinStarExorcistVolume1 = `<article class="g-3 g-3--md mar-b-lg bg-white color-off-black type-sm type-rg--lg">
<figure class="ar-square">
		<a href="javascript:void('Favorite')" class="heart-btn o_requires-login color-mid-gray hover-red z-2 pad-x-md pad-y-sm pos-a t-0 l-0" data-target-id="pr_11837" data-target-title="Twin Star Exorcists, Vol. 1" aria-label="Add Twin Star Exorcists, Vol. 1 to favorites">
	<i aria-hidden="true" class="icon-like"></i>
	<span class="o_votes-up disp-ib color-mid-gray mar-l-sm v-mid type-sm">
		+72
	</span>
</a>


		<span class="product-tag pad-x-rg pad-y-sm pos-a t-0 r-0 z-2 bg-off-black color-white">Series Debut!</span>
	<a tabindex="-1" role="presentation" href="/manga-books/manga/twin-star-exorcists-volume-1/product/3702" class="product-thumb ar-inner type-center">
		<img class="" alt="" src="https://dw9to29mmj727.cloudfront.net/products/1421581744.jpg" style="display: inline-block;">
	</a>
</figure>
<div class="pad-x-md pad-x-lg--lg pad-b-md pad-b-lg--lg">
	<div class="mar-b-sm"><a class="color-mid-gray hover-red">Manga</a></div>
	<a class="color-off-black hover-red" href="/manga-books/manga/twin-star-exorcists-volume-1/product/3702">Twin Star Exorcists, Vol. 1</a>
</div>
</article>`
	twinStarExorcistVolume30 = `<article class="g-3 g-3--md mar-b-lg g-omega bg-white color-off-black type-sm type-rg--lg">
<figure class="ar-square">
		<a href="javascript:void('Favorite')" class="heart-btn o_requires-login color-mid-gray hover-red z-2 pad-x-md pad-y-sm pos-a t-0 l-0" data-target-id="pr_15766" data-target-title="Twin Star Exorcists, Vol. 30" aria-label="Add Twin Star Exorcists, Vol. 30 to favorites">
	<i aria-hidden="true" class="icon-like"></i>
	<span class="o_votes-up disp-ib color-mid-gray mar-l-sm v-mid type-sm">
		+8
	</span>
</a>


		<span class="product-tag pad-x-rg pad-y-sm pos-a t-0 r-0 z-2 bg-mid-gray color-white">Pre-Order</span>
	<a tabindex="-1" role="presentation" href="/manga-books/manga/twin-star-exorcists-volume-30/product/7744" class="product-thumb ar-inner type-center">
		<img class="" alt="" src="https://dw9to29mmj727.cloudfront.net/products/197474311X.jpg" style="display: inline-block;">
	</a>
</figure>
<div class="pad-x-md pad-x-lg--lg pad-b-md pad-b-lg--lg">
	<div class="mar-b-sm"><a class="color-mid-gray hover-red">Manga</a></div>
	<a class="color-off-black hover-red" href="/manga-books/manga/twin-star-exorcists-volume-30/product/7744">Twin Star Exorcists, Vol. 30</a>
</div>
</article>`
	noTitle = `<article class="g-3 g-3--md mar-b-lg bg-white color-off-black type-sm type-rg--lg">
<figure class="ar-square">
		<a href="javascript:void('Favorite')" class="heart-btn o_requires-login color-mid-gray hover-red z-2 pad-x-md pad-y-sm pos-a t-0 l-0" data-target-id="pr_11837" data-target-title="Twin Star Exorcists, Vol. 1" aria-label="Add Twin Star Exorcists, Vol. 1 to favorites">
	<i aria-hidden="true" class="icon-like"></i>
	<span class="o_votes-up disp-ib color-mid-gray mar-l-sm v-mid type-sm">
		+72
	</span>
</a>


		<span class="product-tag pad-x-rg pad-y-sm pos-a t-0 r-0 z-2 bg-off-black color-white">Series Debut!</span>
	<a tabindex="-1" role="presentation" href="/manga-books/manga/twin-star-exorcists-volume-1/product/3702" class="product-thumb ar-inner type-center">
		<img class="" alt="" src="https://dw9to29mmj727.cloudfront.net/products/1421581744.jpg" style="display: inline-block;">
	</a>
</figure>
<div class="pad-x-md pad-x-lg--lg pad-b-md pad-b-lg--lg">
	<div class="mar-b-sm"><a class="color-mid-gray hover-red">Manga</a></div>
	<a class="color-off-black hover-red" href="/manga-books/manga/twin-star-exorcists-volume-1/product/3702"></a>
</div>
</article>`
	noRedirectLink = `<article class="g-3 g-3--md mar-b-lg bg-white color-off-black type-sm type-rg--lg">
<figure class="ar-square">
		<a href="javascript:void('Favorite')" class="heart-btn o_requires-login color-mid-gray hover-red z-2 pad-x-md pad-y-sm pos-a t-0 l-0" data-target-id="pr_11837" data-target-title="Twin Star Exorcists, Vol. 1" aria-label="Add Twin Star Exorcists, Vol. 1 to favorites">
	<i aria-hidden="true" class="icon-like"></i>
	<span class="o_votes-up disp-ib color-mid-gray mar-l-sm v-mid type-sm">
		+72
	</span>
</a>


		<span class="product-tag pad-x-rg pad-y-sm pos-a t-0 r-0 z-2 bg-off-black color-white">Series Debut!</span>
	<a tabindex="-1" role="presentation" href="/manga-books/manga/twin-star-exorcists-volume-1/product/3702" class="product-thumb ar-inner type-center">
		<img class="" alt="" src="https://dw9to29mmj727.cloudfront.net/products/1421581744.jpg" style="display: inline-block;">
	</a>
</figure>
<div class="pad-x-md pad-x-lg--lg pad-b-md pad-b-lg--lg">
	<div class="mar-b-sm"><a class="color-mid-gray hover-red">Manga</a></div>
	<a class="color-off-black hover-red">Twin Star Exorcists, Vol. 1</a>
</div>
</article>`
)

var ParseVolumeHtmlTestCases = map[string]ParseVolumeHtmlTestCase{
	"a volume with no pre-order info should properly have its name and reading link parsed out": {
		InputHtml:            twinStarExorcistVolume1,
		InputSeriesName:      "Twin Star Exorcists",
		ExpectedName:         "Twin Star Exorcists, Vol. 1",
		ExpectedIsReleased:   true,
		ExpectedRedirectLink: "/manga-books/manga/twin-star-exorcists-volume-1/product/3702",
	},
	"a volume with pre-order info should properly have its name and reading link parsed out": {
		InputHtml:            twinStarExorcistVolume30,
		InputSeriesName:      "Twin Star Exorcists",
		ExpectedName:         "Twin Star Exorcists, Vol. 30",
		ExpectedIsReleased:   false,
		ExpectedRedirectLink: "/manga-books/manga/twin-star-exorcists-volume-30/product/7744",
	},
	"a volume with no title results in an error": {
		InputHtml:   noTitle,
		ExpectError: true,
	},
	"a volume with no href results in an error": {
		InputHtml:   noRedirectLink,
		ExpectError: true,
	},
}

func TestParseVolumeHtml(t *testing.T) {
	for name, args := range ParseVolumeHtmlTestCases {
		t.Run(name, func(t *testing.T) {
			actualName, actualRedirectLink, actualIsReleased, err := vizmedia.ParseVolumeHtml(args.InputHtml, args.InputSeriesName, 0)

			assert.Equal(t, err != nil, args.ExpectError)
			if err != nil && !args.ExpectError {
				assert.Fail(t, "An error should not have been thrown: "+err.Error())
			}

			assert.Equal(t, args.ExpectedName, actualName)
			assert.Equal(t, args.ExpectedRedirectLink, actualRedirectLink)
			assert.Equal(t, args.ExpectedIsReleased, actualIsReleased)
		})
	}
}

#!/bin/bash
# Helper script to download images from thebuildmaestro.com
# Since Wix blocks direct downloads, you'll need to:
# 1. Open https://www.thebuildmaestro.com in your browser
# 2. Open Developer Tools (F12) > Network tab
# 3. Filter by "Img" 
# 4. Right-click each image > Copy > Copy image address
# 5. Run: wget -O static/images/hero/devops6.png "<pasted-url>"

echo "Since Wix blocks direct downloads, please download images manually:"
echo ""
echo "1. Open https://www.thebuildmaestro.com in your browser"
echo "2. Open Developer Tools (F12) > Network tab > Filter by 'Img'"
echo "3. Reload the page"
echo "4. Find each image in the network tab, right-click > Copy > Copy image address"
echo "5. Download using: wget -O static/images/path/filename.ext '<pasted-url>'"
echo ""
echo "Required images:"
echo "  - static/images/hero/devops6.png"
echo "  - static/images/hero/me232232.jpg"
echo "  - static/images/companies/gepredix.jpg"
echo "  - static/images/companies/magic-of-macys.jpeg"
echo "  - static/images/companies/expedia400_5.png"
echo "  - static/images/companies/good2.png"
echo "  - static/images/companies/seagullscientific.jpg"
echo "  - static/images/companies/harbor2.jpg"
echo "  - static/images/leanlabs/bbs5.png"
echo "  - static/images/leanlabs/me_working3.jpg"
echo "  - static/images/leanlabs/imagesearch.jpg"
echo "  - static/images/leanlabs/leanlabs1.jpeg"
echo "  - static/images/leanlabs/1483407_850869474975784_6478748118303499751_n.jpg"


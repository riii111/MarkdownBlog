<template>
    <footer class="bg-gray-900 text-white">
        <div class="max-w-7xl mx-auto px-4 pt-16 pb-8">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
                <div class="space-y-6 flex flex-col">
                    <h2 class="text-2xl font-bold">MD Writer</h2>
                    <p class="text-gray-300 max-w-md">
                        The community for developers to share knowledge, best practices,
                        and insights about modern web development.
                    </p>
                </div>
                <div class="lg:flex lg:justify-center">
                    <div class="space-y-4">
                        <h3 class="text-lg font-semibold">Get Updates</h3>
                        <p class="text-gray-400">
                            Subscribe to our newsletter to get the latest articles, tools, and resources.
                        </p>
                        <div class="flex">
                            <UInput type="email" placeholder="Enter your email" variant="outline" class="flex-1" />
                            <CommonButton icon="i-lucide-send" type="submit" label="Subscribe" />
                        </div>
                        <div class="flex justify-between w-full mt-4">
                            <a href="#" class="hover:text-blue-400">
                                <Icon name="lucide:twitter" class="h-5 w-5" />
                            </a>
                            <a href="#" class="hover:text-blue-600">
                                <Icon name="lucide:facebook" class="h-5 w-5" />
                            </a>
                            <a href="#" class="hover:text-pink-500">
                                <Icon name="lucide:instagram" class="h-5 w-5" />
                            </a>
                            <a href="#" class="hover:text-blue-500">
                                <Icon name="lucide:linkedin" class="h-5 w-5" />
                            </a>
                            <a href="#" class="hover:text-gray-400">
                                <Icon name="lucide:github" class="h-5 w-5" />
                            </a>
                        </div>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-8 py-8 border-t border-gray-800">
                <div>
                    <h4 class="font-semibold mb-4">About</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#" class="hover:text-white">About Us</a></li>
                        <li><a href="#" class="hover:text-white">Blog</a></li>
                        <li><a href="#" class="hover:text-white">Team</a></li>
                        <li><a href="#" class="hover:text-white">Careers</a></li>
                    </ul>
                </div>
                <div>
                    <h4 class="font-semibold mb-4">Resources</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#" class="hover:text-white">Documentation</a></li>
                        <li><a href="#" class="hover:text-white">Help Center</a></li>
                        <li><a href="#" class="hover:text-white">Release Notes</a></li>
                        <li><a href="#" class="hover:text-white">Status</a></li>
                    </ul>
                </div>
                <div>
                    <h4 class="font-semibold mb-4">Features</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#" class="hover:text-white">Markdown</a></li>
                        <li><a href="#" class="hover:text-white">Snippets</a></li>
                        <li><a href="#" class="hover:text-white">Integrations</a></li>
                        <li><a href="#" class="hover:text-white">Templates</a></li>
                    </ul>
                </div>
                <div>
                    <h4 class="font-semibold mb-4">Company</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#" class="hover:text-white">Privacy Policy</a></li>
                        <li><a href="#" class="hover:text-white">Terms of Service</a></li>
                        <li><a href="#" class="hover:text-white">Cookie Policy</a></li>
                        <li><a href="#" class="hover:text-white">Contact Us</a></li>
                    </ul>
                </div>
            </div>

        </div>
        <div
            class="border-t border-gray-800 flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <p class="text-gray-400 text-sm">
                © 2024 MD Writer. All rights reserved.
            </p>
            <div class="flex items-center space-x-4 text-sm text-gray-400">
                <a href="#" class="hover:text-white">Terms</a>
                <span>•</span>
                <a href="#" class="hover:text-white">Privacy</a>
                <span>•</span>
                <a href="#" class="hover:text-white">Cookies</a>
            </div>
        </div>
    </footer>
</template>
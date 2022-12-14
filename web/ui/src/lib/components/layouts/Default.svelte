<div class="min-h-screen flex flex-col">
    <div class='bg-blue-700 px-8 flex items-center justify-between py-4 shadow-sm text-white'> App Bar </div>
    <div class='relative flex flex-grow'>
        <nav class='bg-white shadow-sm p-6 space-y-6 w-64'> <slot name="nav"></slot></nav>
        <main class='bg-gray-100 flex-1 p-6'><slot name="main"></slot></main>
    </div>
    <footer class="mt-auto">
        My Footer
    </footer>
</div>